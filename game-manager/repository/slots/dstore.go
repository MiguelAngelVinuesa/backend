package slots

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/goccy/go-json"

	models2 "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/models"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/state/slots"
)

// dstore represents a game bet&state manager utilizing D-Store as the backing store.
// It should be created and initialized only once during app initialization.
type dstore struct {
	logReqResp         bool
	client             *http.Client
	path               string
	apikey             string
	contentType        string
	roundURI           string
	complexInitURI     string
	complexBetURI      string
	complexWinURI      string
	complexCompleteURI string
	roundInitURI       string
	roundCompleteURI   string
	roundNextURI       string
	roundStateURI      string
	sessionStateURI    string
	gamePrefsURI       string
	playerPrefsURI     string
	logger             log.Logger
}

// NewDStore instantiates a new RoundManager which uses D-store as the backing store.
// The uri must be a valid URL without a path; e.g. https://ds.dev.topgaming.team.
// The function will panic if the uri is empty.
// The given timeout must be a reasonable timeout for D-store API calls.
func NewDStore(uri, apikey string, timeout time.Duration, transport http.RoundTripper, logger log.Logger, reqResp bool) slots.RoundManager {
	if uri == "" {
		panic(consts.MsgDsUriMustBeFilled)
	}

	u, _ := url.Parse(uri)
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	prefix := u.Scheme + "://" + u.Host
	if u.Path != "" {
		prefix += u.Path
		if l := len(prefix) - 1; prefix[l] == '/' {
			prefix = prefix[:l]
		}
	}

	return &dstore{
		path:               uri,
		apikey:             apikey,
		client:             &http.Client{Transport: transport, Timeout: timeout},
		contentType:        consts.DefaultContentType,
		roundURI:           prefix + consts.DsRoundURI,
		complexInitURI:     prefix + consts.DsComplexInitURI,
		complexBetURI:      prefix + consts.DsComplexBetURI,
		complexWinURI:      prefix + consts.DsComplexWinURI,
		complexCompleteURI: prefix + consts.DsComplexCompleteURI,
		roundInitURI:       prefix + consts.DsRoundInitURI,
		roundCompleteURI:   prefix + consts.DsRoundCompleteURI,
		roundNextURI:       prefix + consts.DsRoundNextURI,
		roundStateURI:      prefix + consts.DsRoundStateURI,
		sessionStateURI:    prefix + consts.DsSessionStateURI,
		gamePrefsURI:       prefix + consts.DsGameStateURI,
		playerPrefsURI:     prefix + consts.DsPlayerStateURI,
		logger:             logger,
		logReqResp:         reqResp,
	}
}

// PostRound posts a round to D-store, checks the results, and returns the new balance if the round was accepted.
// If the API call fails the function will return false.
func (m *dstore) PostRound(r *slots.Round, debug bool) (string, int64, error) {
	enc, err := models2.MarshallRoundRequest(r, debug)
	if err != nil {
		return m.roundFailed(consts.MsgDsRoundFailed, enc, nil, err)
	}
	return m.postRound(consts.MsgDsRoundFailed, m.roundURI, enc)
}

// PostInitRound posts the initial results of a round to D-store, and returns the new balance if the round was accepted.
// If the API call fails the function will return false.
func (m *dstore) PostInitRound(r *slots.Round, debug bool) (string, int64, error) {
	haveBet := r.TotalBet() > 0
	haveWin := r.TotalWin() > 0

	roundID, balance, err := m.postComplexInit(r, debug, !haveBet && !haveWin)

	if err == nil && haveBet {
		_, balance, err = m.postComplexBet(roundID, r, debug, !haveWin)
	}

	if err == nil && haveWin {
		_, balance, err = m.postComplexWin(roundID, r, debug, true)
	}
	return roundID, balance, err
}

// PostCompleteRound posts the final results of a round to D-store, checks the results, and returns the new balance if the round was accepted.
// If the API call fails the function will return false.
func (m *dstore) PostCompleteRound(r *slots.Round, rs *slots.RoundState, debug bool) (string, int64, error) {
	if rs != nil {
		rs.ResumePlay(len(r.Results()) + 1)
	}

	haveWin := r.TotalWin() > 0
	if haveWin {
		if _, _, err := m.postComplexWin(r.RoundID(), r, debug, true); err != nil {
			return "", 0, err
		}
	}
	return m.postComplexComplete(r, rs, debug, !haveWin)
}

// postComplexInit initializes a new round in D-store, and returns the unqiue round id or an error.
// If the API call fails the function will return false.
func (m *dstore) postComplexInit(r *slots.Round, debug, withData bool) (string, int64, error) {
	enc, err := models2.MarshallComplexInitRequest(r, debug, withData)
	if err != nil {
		return m.roundFailed(consts.MsgDsComplexInitFailed, enc, nil, err)
	}
	return m.postRound(consts.MsgDsComplexInitFailed, m.complexInitURI, enc)
}

// postComplexBet posts a bet for a complex round in D-store, and returns the new balance or an error.
// If the API call fails the function will return false.
func (m *dstore) postComplexBet(roundID string, r *slots.Round, debug, withData bool) (string, int64, error) {
	enc, err := models2.MarshallComplexBetRequest(roundID, r, debug, withData)
	if err != nil {
		return m.roundFailed(consts.MsgDsComplexBetFailed, enc, nil, err)
	}
	return m.postRound(consts.MsgDsComplexBetFailed, m.complexBetURI, enc)
}

// postComplexWin posts a win for a complex round to D-store, and returns the new balance or an error.
// If the API call fails the function will return false.
func (m *dstore) postComplexWin(roundID string, r *slots.Round, debug, withData bool) (string, int64, error) {
	enc, err := models2.MarshallComplexWinRequest(roundID, r, debug, withData)
	if err != nil {
		return m.roundFailed(consts.MsgDsComplexWinFailed, enc, nil, err)
	}
	return m.postRound(consts.MsgDsComplexWinFailed, m.complexWinURI, enc)
}

// postComplexComplete completes a complex round in D-store, and returns the new balance or an error.
// If the API call fails the function will return false.
func (m *dstore) postComplexComplete(r *slots.Round, rs *slots.RoundState, debug, withData bool) (string, int64, error) {
	enc, err := models2.MarshallComplexCompleteRequest(r, rs, debug, withData)
	if err != nil {
		return m.roundFailed(consts.MsgDsComplexCompleteFailed, enc, nil, err)
	}
	return m.postRound(consts.MsgDsComplexCompleteFailed, m.complexCompleteURI, enc)
}

func (m *dstore) postRound(msg, uri string, enc *zjson.Encoder) (string, int64, error) {
	defer enc.Release()

	req, err := m.newRequest(http.MethodPost, uri, bytes.NewReader(enc.Bytes()))
	if err != nil {
		return m.roundFailed(msg, enc, nil, err)
	}

	resp, err2 := m.httpRequest(req)
	if err2 != nil || resp == nil {
		return m.roundFailed(msg, enc, resp, err2)
	}
	defer resp.Body.Close()

	roundID, balance, err3 := models2.UnmarshallRoundResponse(resp)
	if err3 != nil {
		return m.roundFailed(msg, enc, nil, err3)
	}

	return roundID, balance, nil
}

// PostRoundNext moves to the next spin of a multi-spin result.
// If the API call fails the function will return false.
func (m *dstore) PostRoundNext(sessionID, roundID string, roundState *slots.RoundState, spinSeq int) (*slots.RoundResult, int64, error) {
	enc, err := models2.MarshallRoundNextRequest(sessionID, roundID, roundState, spinSeq)
	defer enc.Release()
	if err != nil {
		return m.roundNextFailed(enc, nil, err)
	}

	req, err2 := m.newRequest(http.MethodPost, m.roundNextURI, bytes.NewReader(enc.Bytes()))
	if err2 != nil {
		return m.roundNextFailed(enc, nil, err2)
	}

	resp, err3 := m.httpRequest(req)
	if err3 != nil || resp == nil {
		return m.roundNextFailed(enc, resp, err3)
	}
	defer resp.Body.Close()

	result, balance, err4 := models2.UnmarshallRoundNextResponse(resp)
	if err4 != nil {
		return m.roundNextFailed(enc, nil, err4)
	}

	return result, balance, nil
}

// GetRoundState retrieves the game state for the current session from D-Store.
// If the API call fails the function will return false.
func (m *dstore) GetRoundState(sessionID, roundID string) (*slots.RoundState, error) {
	uri := fmt.Sprintf("%s?session=%s&round=%s", m.roundStateURI, sessionID, roundID)
	req, err := m.newRequest(http.MethodGet, uri, nil)
	if err != nil {
		return m.getRoundStateFailed(sessionID, roundID, nil, err)
	}

	resp, err2 := m.httpRequest(req)
	if err2 != nil || resp == nil {
		return m.getRoundStateFailed(sessionID, roundID, resp, err2)
	}
	defer resp.Body.Close()

	state, err3 := models2.UnmarshallGetRoundStateResponse(sessionID, roundID, resp)
	if err3 != nil {
		return m.getRoundStateFailed(sessionID, roundID, nil, err3)
	}

	return state, nil
}

// PutGameState updates the game state for the current session in D-Store.
// If the API call fails the function will return false.
func (m *dstore) PutGameState(sessionID string, state *slots.GameState) error {
	enc, err := models2.MarshallSessionStateRequest(sessionID, state)
	defer enc.Release()
	if err != nil {
		return m.putGameStateFailed(enc, nil, err)
	}

	req, err2 := m.newRequest(http.MethodPut, m.sessionStateURI, bytes.NewReader(enc.Bytes()))
	if err2 != nil {
		return m.putGameStateFailed(enc, nil, err2)
	}

	resp, err3 := m.httpRequest(req)
	if err != nil || resp == nil {
		return m.putGameStateFailed(enc, resp, err3)
	}
	defer resp.Body.Close()

	if _, err = models2.UnmarshallPutSessionStateResponse(resp); err != nil {
		return m.putGameStateFailed(enc, nil, err)
	}
	return nil
}

// GetGameState retrieves the game state for the current session from D-Store.
// If the API call fails the function will return false.
func (m *dstore) GetGameState(sessionID string) (*slots.GameState, error) {
	uri := fmt.Sprintf("%s?session=%s", m.sessionStateURI, sessionID)
	req, err := m.newRequest(http.MethodGet, uri, nil)
	if err != nil {
		return m.getGameStateFailed(sessionID, nil, err)
	}

	resp, err2 := m.httpRequest(req)
	if err2 != nil || resp == nil {
		if resp != nil && resp.StatusCode == http.StatusInternalServerError {
			return models2.EmptySessionState(), nil
		}
		return m.getGameStateFailed(sessionID, resp, err2)
	}
	defer resp.Body.Close()

	state, err3 := models2.UnmarshallGetSessionStateResponse(resp)
	if err3 != nil {
		return m.getGameStateFailed(sessionID, nil, err3)
	}
	return state, nil
}

// PutGamePrefs stores the player game preferences in D-Store.
// If the API call fails the function will return false.
func (m *dstore) PutGamePrefs(sessionID string, state *slots.GamePrefs) error {
	enc, err := models2.MarshallGamePrefsRequest(sessionID, state)
	defer enc.Release()
	if err != nil {
		return m.putGamePrefsFailed(enc, nil, err)
	}

	req, err2 := m.newRequest(http.MethodPut, m.gamePrefsURI, bytes.NewReader(enc.Bytes()))
	if err2 != nil {
		return m.putGamePrefsFailed(enc, nil, err2)
	}

	resp, err3 := m.httpRequest(req)
	if err3 != nil || resp == nil {
		return m.putGamePrefsFailed(enc, resp, err3)
	}
	defer resp.Body.Close()

	if _, err = models2.UnmarshallPutGamePrefsResponse(resp); err != nil {
		return m.putGamePrefsFailed(enc, nil, err)
	}
	return nil
}

// GetGamePrefs retrieves the player game preferences from D-Store.
// If the API call fails the function will return false.
func (m *dstore) GetGamePrefs(sessionID string) (string, string, *slots.GamePrefs, error) {
	uri := fmt.Sprintf("%s?session=%s", m.gamePrefsURI, sessionID)
	req, err := m.newRequest(http.MethodGet, uri, nil)
	if err != nil {
		return m.getGamePrefsFailed(sessionID, nil, err)
	}

	resp, err2 := m.httpRequest(req)
	if err2 != nil || resp == nil {
		if resp != nil && resp.StatusCode == http.StatusInternalServerError {
			return "", "", models2.EmptyGamePrefs(), nil
		}
		return m.getGamePrefsFailed(sessionID, resp, err2)
	}
	defer resp.Body.Close()

	casino, player, state, err3 := models2.UnmarshallGetGamePrefsResponse(resp)
	if err3 != nil {
		return m.getGamePrefsFailed(sessionID, nil, err3)
	}
	return casino, player, state, nil
}

// PutPlayerPrefs stores the player global preferences in D-Store.
// If the API call fails the function will return false.
func (m *dstore) PutPlayerPrefs(sessionID string, state map[string]string) error {
	enc, err := models2.MarshallPlayerPrefsRequest(sessionID, state)
	defer enc.Release()
	if err != nil {
		return m.putPlayerPrefsFailed(enc, nil, err)
	}

	req, err2 := m.newRequest(http.MethodPut, m.playerPrefsURI, bytes.NewReader(enc.Bytes()))
	if err2 != nil {
		return m.putPlayerPrefsFailed(enc, nil, err2)
	}

	resp, err3 := m.httpRequest(req)
	if err3 != nil || resp == nil {
		return m.putPlayerPrefsFailed(enc, resp, err3)
	}
	defer resp.Body.Close()

	if _, err = models2.UnmarshallPutPlayerPrefsResponse(resp); err != nil {
		return m.putPlayerPrefsFailed(enc, nil, err)
	}
	return nil
}

// GetPlayerPrefs retrieves the player global preferences from D-Store.
// If the API call fails the function will return false.
func (m *dstore) GetPlayerPrefs(sessionID string) (map[string]string, error) {
	uri := fmt.Sprintf("%s?session=%s", m.playerPrefsURI, sessionID)
	req, err := m.newRequest(http.MethodGet, uri, nil)
	if err != nil {
		return m.getPlayerPrefsFailed(sessionID, nil, err)
	}

	resp, err2 := m.httpRequest(req)
	if err2 != nil || resp == nil {
		if resp != nil && resp.StatusCode == http.StatusInternalServerError {
			return map[string]string{}, nil
		}
		return m.getPlayerPrefsFailed(sessionID, resp, err2)
	}
	defer resp.Body.Close()

	state, err3 := models2.UnmarshallGetPlayerPrefsResponse(resp)
	if err3 != nil {
		return m.getPlayerPrefsFailed(sessionID, nil, err3)
	}
	return state, nil
}

func (m *dstore) newRequest(method string, uri string, body *bytes.Reader) (*http.Request, error) {
	var req *http.Request
	var err error

	haveBody := body != nil && body.Len() > 0
	if haveBody {
		req, err = http.NewRequest(method, uri, body)
	} else {
		req, err = http.NewRequest(method, uri, nil)
	}
	if err != nil {
		return req, err
	}

	if haveBody {
		req.Header.Add("Content-Type", m.contentType)
	}

	req.Header.Add("Accept-Encoding", m.contentType)
	req.Header.Add("X-API-Key", m.apikey)

	return req, nil
}

func (m *dstore) httpRequest(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	backoff := 10 * time.Millisecond

	for {
		resp, err = m.client.Do(req)
		if err != nil || resp == nil {
			return resp, err
		}

		if resp.StatusCode == http.StatusOK {
			break
		}

		switch resp.StatusCode {
		case http.StatusBadGateway:
			// retry
		default:
			return resp, fmt.Errorf(consts.MsgDsInvalidStatus, resp.StatusCode)
		}

		time.Sleep(backoff)
		if backoff < 2*time.Second {
			backoff *= 2
		}
	}

	return resp, nil
}

func (m *dstore) roundFailed(label string, enc *zjson.Encoder, resp *http.Response, err error) (string, int64, error) {
	err = m.errorFromResponse(err, resp)
	if m.logger != nil {
		m.logger.Error(label, consts.FieldRequest, string(enc.Bytes()), consts.FieldError, err)
	}
	return "", 0, err
}

func (m *dstore) roundNextFailed(enc *zjson.Encoder, resp *http.Response, err error) (*slots.RoundResult, int64, error) {
	err = m.errorFromResponse(err, resp)
	if m.logger != nil {
		m.logger.Error(consts.MsgDsRoundNextFailed, consts.FieldRequest, string(enc.Bytes()), consts.FieldError, err)
	}
	return nil, 0, err
}

func (m *dstore) getRoundStateFailed(sessionID, roundID string, resp *http.Response, err error) (*slots.RoundState, error) {
	err = m.errorFromResponse(err, resp)
	if m.logger != nil {
		m.logger.Error(consts.MsgDsGetRoundStateFailed, consts.FieldSession, sessionID, consts.FieldRound, roundID, consts.FieldError, err)
	}
	return nil, err
}

func (m *dstore) putGameStateFailed(enc *zjson.Encoder, resp *http.Response, err error) error {
	err = m.errorFromResponse(err, resp)
	if m.logger != nil {
		m.logger.Error(consts.MsgDsPutSessionStateFailed, consts.FieldRequest, string(enc.Bytes()), consts.FieldError, err)
	}
	return err
}

func (m *dstore) getGameStateFailed(sessionID string, resp *http.Response, err error) (*slots.GameState, error) {
	err = m.errorFromResponse(err, resp)
	if m.logger != nil {
		m.logger.Error(consts.MsgDsGetSessionStateFailed, consts.FieldSession, sessionID, consts.FieldError, err)
	}
	return nil, err
}

func (m *dstore) putGamePrefsFailed(enc *zjson.Encoder, resp *http.Response, err error) error {
	err = m.errorFromResponse(err, resp)
	if m.logger != nil {
		m.logger.Error(consts.MsgDsPutGamePrefsFailed, consts.FieldRequest, string(enc.Bytes()), consts.FieldError, err)
	}
	return err
}

func (m *dstore) getGamePrefsFailed(sessionID string, resp *http.Response, err error) (string, string, *slots.GamePrefs, error) {
	err = m.errorFromResponse(err, resp)
	if m.logger != nil {
		m.logger.Error(consts.MsgDsGetGamePrefsFailed, consts.FieldSession, sessionID, consts.FieldError, err)
	}
	return "", "", nil, err
}

func (m *dstore) putPlayerPrefsFailed(enc *zjson.Encoder, resp *http.Response, err error) error {
	err = m.errorFromResponse(err, resp)
	if m.logger != nil {
		m.logger.Error(consts.MsgDsPutPlayerPrefsFailed, consts.FieldRequest, string(enc.Bytes()), consts.FieldError, err)
	}
	return err
}

func (m *dstore) getPlayerPrefsFailed(sessionID string, resp *http.Response, err error) (map[string]string, error) {
	err = m.errorFromResponse(err, resp)
	if m.logger != nil {
		m.logger.Error(consts.MsgDsGetPlayerPrefsFailed, consts.FieldSession, sessionID, consts.FieldError, err)
	}
	return nil, err
}

func (m *dstore) errorFromResponse(err error, resp *http.Response) error {
	out := &slots.APIerror{Err: err, Level: "F"}

	if resp != nil {
		out.Status = resp.StatusCode

		b, err2 := io.ReadAll(resp.Body)
		if err2 != nil {
			return out
		}

		data := make(map[string]any, 10)
		if err2 = json.Unmarshal(b, &data); err2 != nil {
			return out
		}

		if s, ok := data["error"].(string); ok {
			out.Message = s
		} else if s, ok = data["message"].(string); ok {
			out.Message = s
		} else if s, ok = data["errorMessage"].(string); ok {
			out.Message = s
		}

		if f, ok := data["error-code"].(float64); ok {
			out.Code = int(f)
		} else if f, ok = data["errorCode"].(float64); ok {
			out.Code = int(f)
		}

		if s, ok := data["error-level"].(string); ok {
			out.Level = s
		} else if s, ok = data["errorLevel"].(string); ok {
			out.Level = s
		}
	}

	if out.Level == "F" {
		switch out.Code {
		case 1211, 1212, 1213, 1221, 1232: // D-store retryable.
			out.Level = "R"
		}
	}

	return out
}
