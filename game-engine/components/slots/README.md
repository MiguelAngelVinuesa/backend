# SLOTS #

This directory contains components to build fluid slot machine backends.


## Terminology / naming conventions ##

| Term           | Description                                                                                      |
|----------------|--------------------------------------------------------------------------------------------------|
| ActionKind     | kind of action; e.g. payout, free spin, bonus game                                               |
| ActionStage    | stage of the spin processing at which an action will be tested                                   |
| AllPaylines    | a feature that enables all possible paylines; on a 5x3 grid this is 243 possible paylines        |
| AllScatter     | this feature turns all symbols into scatter symbols; e.g. no paylines, just scatter payouts      |
| Bomb           | a special symbol that can morph the grid tiles around it                                         |
| BonusGame      | an award from one of the triggers that allows to play a bonus game like "double or nothing"      |
| BonusSymbol    | a randomly selected symbol for free spins that acts as a special scatter                         |
| BothWays       | indicates that paylines may be calculated both LTR and RTL                                       |
| CascadingReels | a feature where after a spin symbols on the grid are removed to make way for new symbols         |
| ClearAction    | an action that when triggered clears part of the spin result                                     |
| Direction      | direction from which the winning symbols on a payline are counted (LTR, RTL, BothWays)           |
| DoubleSpin     | a feature where the spin is performed in two steps (such as in ChaCha Bomb)                      |
| ExpandAfter    | wild expansion happens after the paylines have been matched                                      |
| ExpandBefore   | wild expansion happens before the paylines have been matched                                     |
| ExpandingWilds | a mechanism where wild symbols are expanded across the reel they appear in                       |
| FreeSpin       | an additional free-of-charge spin awarded from one of the triggers                               |
| Hero           | a symbol which if displayed may cause wild symbols to be expanded                                |
| HeroScatter    | a symbol that acts as both a hero and a scatter symbol                                           |
| Hot            | indicates a reel as hot, meaning it will have special powers during the following spins          |
| HighestPayout  | a feature that awards highest payout when substituting wilds for the highest paying symbol       |
| JumpingWild    | a feature where existing wild symbol(s) "jumps" around the grid on each spin                     |
| Locked         | flag to indicate a reel is locked and won't "spin"                                               |
| MorphAction    | an action that can morph symbols on the grid before any other actions are tested                 |
| Multiplier     | a factor used to increase one or more payouts                                                    |
| NoRepeat       | number to indicate how many randomly selected symbols on a reel must not repeat consecutively    |
| PaidAction     | an action that is activated by the user                                                          |
| Payline        | virtual line going through the symbols on the reels that represents a potential payout           |
| PayLTR         | left-to-right direction for paylines                                                             |
| PayRTL         | right-to-left direction for paylines                                                             |
| Payout         | awarded payout from a payline or trigger                                                         |
| Reel           | a reel with symbols in a slot machine                                                            |
| ReelCount      | the number of reels in a slot machine                                                            |
| RowCount       | the number of rows in a slot machine                                                             |
| Scatter        | a symbol that can trigger a payout or bonus game without being on a specific payline             |
| ScatterAction  | an action that activates when the required amount of scatter symbols appear on the reels         |
| Slots          | slot(s) machine / fruit machine / poker machine / one-armed bandit                               |
| Spin           | result of a "spin" of the unlocked reels                                                         |
| SpinAction     | an action indicates how and when certain game features are triggered and executed                |
| Split          | a symbol that can represent two or more specific symbols, like a limited wild symbol             |
| Standard       | a standard symbol with no special function(s) other than to warrant payouts                      |
| StartGrid      | a pre-calculated grid to start a spin with                                                       |
| Sticky         | flag to indicate that a symbol is "sticky" (e.g. remains in place on the next spin)              |
| StickyAction   | an action that can mkae one or more symbols on a spin result become sticky                       |
| Symbol         | symbol displayed on the reels of a slot machine                                                  |
| SymbolKind     | the kind of symbol, such as Wild, Hero, Scatter                                                  |
| SymbolSet      | set of possible symbols used by a slot machine                                                   |
| Weighting      | weighs for a symbol to determine the frequency of its random occurrence on each reel             |
| Wild           | a symbol that can represent any of the other symbols and therefor always counts towards paylines |
| WildScatter    | a symbol that acts as both a wild and a scatter symbol                                           |
| WildAction     | an action that is activated when one or more new wild symbols appear on the reels                |
| Winline        | a payline with enough consecutive equal symbols including wilds&splits to constitutes a payout   |
