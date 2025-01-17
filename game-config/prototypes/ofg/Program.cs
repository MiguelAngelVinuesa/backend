//using System;
//using System.Collections.Generic;
using System.Diagnostics;
//using System.Linq;
//using System.Threading.Tasks;

namespace Sim_Round
{
    internal class Program
    {
        static void Main(string[] args)
        {

            var seed = new Random(DateTime.Now.Millisecond);

            float[,] result = new float[10, 6];

            Parallel.For(0, 6, ip =>
            {
                Console.WriteLine($"Task {Task.CurrentId} is running on core {ip}");
                // ------- Your task logic goes here ---------------------------------------------------

                int tajm = seed.Next();
                var rnd = new Random(tajm);
                Console.WriteLine(tajm);

                int ntot = 50000000;

                long total = 0;
                long bg_total = 0;
                long fg_total = 0;

                long std_total = 0;

                int plus = 0;
                int cap = 0;

                int nx = 5;
                int ny = 3;

                int bet = 10;

                List<int[]> bgReelSet1 = new List<int[]>()
                {
                    new int[] {4, 5, 1, 7, 9, 8, 6, 1, 7, 8, 3, 6, 4, 0, 9, 5, 1, 6, 8, 3, 9, 5, 1, 6, 4, 2, 8, 7, 3, 4, 5, 2, 6, 7, 3, 8, 9, 0, 7, 4, 3, 9, 5, 1, 6, 8, 2, 4, 5, 1, 6, 7, 3, 5, 9, 1, 8, 9, 2, 6, 4, 0, 7, 8, 2, 5, 4, 1, 9, 7, 3, 6, 9, 2, 4, 5, 3, 6, 8, 9, 7, 2, 8, 6, 1, 5, 4, 3, 8, 7, 1, 9, 8, 2, 7, 9, 0, 6, 7, 3, 4, 5, 1, 9, 8, 2, 4, 7, 1, 5, 8, 2, 9, 6, 3, 1, 5, 9, 3, 7, 6, 1, 4, 8, 1, 2, 9, 4, 1, 6, 8, 2, 7, 5, 1, 3, 9, 8, 3, 7, 4, 1, 5, 6, 1, 2, 4, 9, 3, 8, 6, 2, 5, 7, 0, 1, 8, 9, 1, 7, 4, 1, 6, 5, 2, 3, 4, 7, 1, 5, 8, 2, 9, 6, 1, 3, 5, 9, 3, 7, 6, 1, 4, 8, 2, 1, 9, 4, 1, 6, 8, 2, 7, 5, 3, 9, 8, 3, 7, 4, 1, 5, 6, 2, 4, 9, 3, 8, 6, 2, 5, 7, 1, 0, 8, 9, 1, 7, 4, 1, 6, 5, 3, 2},
                    new int[] {4, 5, 0, 7, 9, 8, 6, 3, 7, 8, 0, 6, 4, 2, 9, 5, 1, 6, 8, 3, 9, 5, 0, 6, 4, 3, 8, 7, 2, 4, 5, 6, 7, 3, 8, 9, 0, 7, 4, 3, 9, 5, 1, 6, 8, 2, 4, 5, 0, 6, 7, 2, 5, 9, 0, 8, 9, 3, 6, 4, 7, 8, 2, 5, 4, 0, 9, 7, 3, 6, 9, 2, 4, 5, 3, 6, 8, 7, 9, 2, 8, 6, 1, 5, 4, 3, 8, 7, 9, 6, 2, 7, 9, 0, 8, 7, 3, 4, 5, 0, 9, 8, 2, 4, 7, 1, 5, 8, 2, 9, 6, 3, 0, 5, 9, 3, 7, 6, 0, 4, 8, 0, 2, 9, 4, 0, 6, 8, 2, 7, 5, 0, 3, 9, 8, 3, 7, 4, 0, 5, 6, 0, 2, 4, 9, 3, 8, 6, 2, 5, 7, 0, 1, 8, 9, 0, 7, 4, 0, 6, 5, 2, 3, 4, 7, 0, 5, 8, 2, 9, 6, 0, 3, 5, 9, 3, 7, 6, 0, 4, 8, 2, 9, 4, 0, 6, 8, 2, 7, 5, 3, 0, 9, 8, 3, 7, 4, 0, 5, 6, 2, 4, 9, 3, 8, 6, 2, 5, 7, 1, 0, 8, 9, 0, 7, 4, 0, 6, 5, 2, 3},
                    new int[] {4, 5, 3, 7, 9, 1, 0, 8, 6, 0, 7, 8, 1, 6, 4, 0, 1, 9, 5, 1, 6, 8, 3, 9, 5, 1, 6, 4, 2, 0, 8, 7, 3, 4, 5, 2, 6, 7, 3, 8, 9, 0, 1, 7, 4, 3, 9, 5, 1, 6, 8, 2, 4, 5, 0, 2, 6, 7, 3, 5, 9, 1, 8, 9, 2, 6, 4, 0, 1, 7, 8, 2, 5, 4, 1, 9, 7, 3, 6, 9, 2, 0, 4, 5, 3, 6, 8, 0, 9, 7, 2, 8, 6, 0, 5, 4, 3, 0, 8, 7, 1, 9, 8, 2, 7, 9, 0, 1, 6, 7, 3, 4, 5, 1, 9, 8, 2, 3, 4, 7, 0, 5, 8, 2, 9, 6, 3, 2, 5, 9, 3, 7, 6, 1, 4, 8, 0, 2, 9, 4, 1, 6, 8, 2, 7, 5, 0, 3, 9, 8, 3, 7, 4, 0, 5, 6, 1, 2, 4, 9, 3, 8, 6, 2, 5, 7, 0, 1, 8, 9, 0, 7, 4, 1, 6, 5, 2, 3, 4, 7, 0, 5, 8, 2, 9, 6, 1, 3, 5, 9, 3, 7, 6, 1, 4, 8, 2, 0, 9, 4, 1, 6, 8, 2, 7, 5, 3, 0, 9, 8, 3, 7, 4, 0, 5, 6, 2, 1, 4, 9, 3, 8, 6, 2, 5, 7, 1, 0, 8, 9, 0, 7, 4, 1, 6, 5, 3, 2},
                    new int[] {4, 5, 3, 7, 9, 1, 0, 8, 6, 0, 7, 8, 1, 6, 4, 0, 1, 9, 5, 1, 6, 8, 3, 9, 5, 1, 6, 4, 2, 0, 8, 7, 3, 4, 5, 2, 6, 7, 3, 8, 9, 0, 1, 7, 4, 3, 9, 5, 1, 6, 8, 2, 4, 5, 0, 2, 6, 7, 3, 5, 9, 1, 8, 9, 2, 6, 4, 0, 1, 7, 8, 2, 5, 4, 1, 9, 7, 3, 6, 9, 2, 0, 4, 5, 3, 6, 8, 0, 9, 7, 2, 8, 6, 0, 5, 4, 3, 0, 8, 7, 1, 9, 8, 2, 7, 9, 0, 1, 6, 7, 3, 4, 5, 1, 9, 8, 2, 3, 4, 7, 0, 5, 8, 2, 9, 6, 3, 2, 5, 9, 3, 7, 6, 1, 4, 8, 0, 2, 9, 4, 1, 6, 8, 2, 7, 5, 0, 3, 9, 8, 3, 7, 4, 0, 5, 6, 1, 2, 4, 9, 3, 8, 6, 2, 5, 7, 0, 1, 8, 9, 0, 7, 4, 1, 6, 5, 2, 3, 4, 7, 0, 5, 8, 2, 9, 6, 1, 3, 5, 9, 3, 7, 6, 1, 4, 8, 2, 0, 9, 4, 1, 6, 8, 2, 7, 5, 3, 0, 9, 8, 3, 7, 4, 0, 5, 6, 2, 1, 4, 9, 3, 8, 6, 2, 5, 7, 1, 0, 8, 9, 0, 7, 4, 1, 6, 5, 3, 2},
                    new int[] {4, 5, 3, 7, 9, 1, 0, 8, 6, 0, 7, 8, 1, 6, 4, 0, 1, 9, 5, 1, 6, 8, 3, 9, 5, 1, 6, 4, 2, 0, 8, 7, 3, 4, 5, 2, 6, 7, 3, 8, 9, 0, 1, 7, 4, 3, 9, 5, 1, 6, 8, 2, 4, 5, 0, 2, 6, 7, 3, 5, 9, 1, 8, 9, 2, 6, 4, 0, 1, 7, 8, 2, 5, 4, 1, 9, 7, 3, 6, 9, 2, 0, 4, 5, 3, 6, 8, 0, 9, 7, 2, 8, 6, 0, 5, 4, 3, 0, 8, 7, 1, 9, 8, 2, 7, 9, 0, 1, 6, 7, 3, 4, 5, 1, 9, 8, 2, 3, 4, 7, 0, 5, 8, 2, 9, 6, 3, 2, 5, 9, 3, 7, 6, 1, 4, 8, 0, 2, 9, 4, 1, 6, 8, 2, 7, 5, 0, 3, 9, 8, 3, 7, 4, 0, 5, 6, 1, 2, 4, 9, 3, 8, 6, 2, 5, 7, 0, 1, 8, 9, 0, 7, 4, 1, 6, 5, 2, 3, 4, 7, 0, 5, 8, 2, 9, 6, 1, 3, 5, 9, 3, 7, 6, 1, 4, 8, 2, 0, 9, 4, 1, 6, 8, 2, 7, 5, 3, 0, 9, 8, 3, 7, 4, 0, 5, 6, 2, 1, 4, 9, 3, 8, 6, 2, 5, 7, 1, 0, 8, 9, 0, 7, 4, 1, 6, 5, 3, 2}
                };

                List<int[]> bgReelSet2 = new List<int[]>()
                {
                    new int[] {4, 5, 0, 7, 9, 8, 6, 3, 7, 8, 0, 6, 4, 2, 9, 5, 1, 6, 8, 3, 9, 5, 0, 6, 4, 3, 8, 7, 2, 4, 5, 6, 7, 3, 8, 9, 0, 7, 4, 3, 9, 5, 1, 6, 8, 2, 4, 5, 0, 6, 7, 2, 5, 9, 0, 8, 9, 3, 6, 4, 7, 8, 2, 5, 4, 0, 9, 7, 3, 6, 9, 2, 4, 5, 3, 6, 8, 7, 9, 2, 8, 6, 1, 5, 4, 3, 8, 7, 9, 6, 2, 7, 9, 0, 8, 7, 3, 4, 5, 0, 9, 8, 2, 4, 7, 1, 5, 8, 2, 9, 6, 3, 0, 5, 9, 3, 7, 6, 0, 4, 8, 0, 2, 9, 4, 0, 6, 8, 2, 7, 5, 0, 3, 9, 8, 3, 7, 4, 0, 5, 6, 0, 2, 4, 9, 3, 8, 6, 2, 5, 7, 0, 1, 8, 9, 0, 7, 4, 0, 6, 5, 2, 3, 4, 7, 0, 5, 8, 2, 9, 6, 0, 3, 5, 9, 3, 7, 6, 0, 4, 8, 2, 9, 4, 0, 6, 8, 2, 7, 5, 3, 0, 9, 8, 3, 7, 4, 0, 5, 6, 2, 4, 9, 3, 8, 6, 2, 5, 7, 1, 0, 8, 9, 0, 7, 4, 0, 6, 5, 2, 3},
                    new int[] {4, 5, 1, 7, 9, 8, 6, 1, 7, 8, 3, 6, 4, 0, 9, 5, 1, 6, 8, 3, 9, 5, 1, 6, 4, 2, 8, 7, 3, 4, 5, 2, 6, 7, 3, 8, 9, 0, 7, 4, 3, 9, 5, 1, 6, 8, 2, 4, 5, 1, 6, 7, 3, 5, 9, 1, 8, 9, 2, 6, 4, 0, 7, 8, 2, 5, 4, 1, 9, 7, 3, 6, 9, 2, 4, 5, 3, 6, 8, 9, 7, 2, 8, 6, 1, 5, 4, 3, 8, 7, 1, 9, 8, 2, 7, 9, 0, 6, 7, 3, 4, 5, 1, 9, 8, 2, 4, 7, 1, 5, 8, 2, 9, 6, 3, 1, 5, 9, 3, 7, 6, 1, 4, 8, 1, 2, 9, 4, 1, 6, 8, 2, 7, 5, 1, 3, 9, 8, 3, 7, 4, 1, 5, 6, 1, 2, 4, 9, 3, 8, 6, 2, 5, 7, 0, 1, 8, 9, 1, 7, 4, 1, 6, 5, 2, 3, 4, 7, 1, 5, 8, 2, 9, 6, 1, 3, 5, 9, 3, 7, 6, 1, 4, 8, 2, 1, 9, 4, 1, 6, 8, 2, 7, 5, 3, 9, 8, 3, 7, 4, 1, 5, 6, 2, 4, 9, 3, 8, 6, 2, 5, 7, 1, 0, 8, 9, 1, 7, 4, 1, 6, 5, 3, 2},
                    new int[] {4, 5, 3, 7, 9, 1, 0, 8, 6, 0, 7, 8, 1, 6, 4, 0, 1, 9, 5, 1, 6, 8, 3, 9, 5, 1, 6, 4, 2, 0, 8, 7, 3, 4, 5, 2, 6, 7, 3, 8, 9, 0, 1, 7, 4, 3, 9, 5, 1, 6, 8, 2, 4, 5, 0, 2, 6, 7, 3, 5, 9, 1, 8, 9, 2, 6, 4, 0, 1, 7, 8, 2, 5, 4, 1, 9, 7, 3, 6, 9, 2, 0, 4, 5, 3, 6, 8, 0, 9, 7, 2, 8, 6, 0, 5, 4, 3, 0, 8, 7, 1, 9, 8, 2, 7, 9, 0, 1, 6, 7, 3, 4, 5, 1, 9, 8, 2, 3, 4, 7, 0, 5, 8, 2, 9, 6, 3, 2, 5, 9, 3, 7, 6, 1, 4, 8, 0, 2, 9, 4, 1, 6, 8, 2, 7, 5, 0, 3, 9, 8, 3, 7, 4, 0, 5, 6, 1, 2, 4, 9, 3, 8, 6, 2, 5, 7, 0, 1, 8, 9, 0, 7, 4, 1, 6, 5, 2, 3, 4, 7, 0, 5, 8, 2, 9, 6, 1, 3, 5, 9, 3, 7, 6, 1, 4, 8, 2, 0, 9, 4, 1, 6, 8, 2, 7, 5, 3, 0, 9, 8, 3, 7, 4, 0, 5, 6, 2, 1, 4, 9, 3, 8, 6, 2, 5, 7, 1, 0, 8, 9, 0, 7, 4, 1, 6, 5, 3, 2},
                    new int[] {4, 5, 3, 7, 9, 1, 0, 8, 6, 0, 7, 8, 1, 6, 4, 0, 1, 9, 5, 1, 6, 8, 3, 9, 5, 1, 6, 4, 2, 0, 8, 7, 3, 4, 5, 2, 6, 7, 3, 8, 9, 0, 1, 7, 4, 3, 9, 5, 1, 6, 8, 2, 4, 5, 0, 2, 6, 7, 3, 5, 9, 1, 8, 9, 2, 6, 4, 0, 1, 7, 8, 2, 5, 4, 1, 9, 7, 3, 6, 9, 2, 0, 4, 5, 3, 6, 8, 0, 9, 7, 2, 8, 6, 0, 5, 4, 3, 0, 8, 7, 1, 9, 8, 2, 7, 9, 0, 1, 6, 7, 3, 4, 5, 1, 9, 8, 2, 3, 4, 7, 0, 5, 8, 2, 9, 6, 3, 2, 5, 9, 3, 7, 6, 1, 4, 8, 0, 2, 9, 4, 1, 6, 8, 2, 7, 5, 0, 3, 9, 8, 3, 7, 4, 0, 5, 6, 1, 2, 4, 9, 3, 8, 6, 2, 5, 7, 0, 1, 8, 9, 0, 7, 4, 1, 6, 5, 2, 3, 4, 7, 0, 5, 8, 2, 9, 6, 1, 3, 5, 9, 3, 7, 6, 1, 4, 8, 2, 0, 9, 4, 1, 6, 8, 2, 7, 5, 3, 0, 9, 8, 3, 7, 4, 0, 5, 6, 2, 1, 4, 9, 3, 8, 6, 2, 5, 7, 1, 0, 8, 9, 0, 7, 4, 1, 6, 5, 3, 2},
                    new int[] {4, 5, 3, 7, 9, 1, 0, 8, 6, 0, 7, 8, 1, 6, 4, 0, 1, 9, 5, 1, 6, 8, 3, 9, 5, 1, 6, 4, 2, 0, 8, 7, 3, 4, 5, 2, 6, 7, 3, 8, 9, 0, 1, 7, 4, 3, 9, 5, 1, 6, 8, 2, 4, 5, 0, 2, 6, 7, 3, 5, 9, 1, 8, 9, 2, 6, 4, 0, 1, 7, 8, 2, 5, 4, 1, 9, 7, 3, 6, 9, 2, 0, 4, 5, 3, 6, 8, 0, 9, 7, 2, 8, 6, 0, 5, 4, 3, 0, 8, 7, 1, 9, 8, 2, 7, 9, 0, 1, 6, 7, 3, 4, 5, 1, 9, 8, 2, 3, 4, 7, 0, 5, 8, 2, 9, 6, 3, 2, 5, 9, 3, 7, 6, 1, 4, 8, 0, 2, 9, 4, 1, 6, 8, 2, 7, 5, 0, 3, 9, 8, 3, 7, 4, 0, 5, 6, 1, 2, 4, 9, 3, 8, 6, 2, 5, 7, 0, 1, 8, 9, 0, 7, 4, 1, 6, 5, 2, 3, 4, 7, 0, 5, 8, 2, 9, 6, 1, 3, 5, 9, 3, 7, 6, 1, 4, 8, 2, 0, 9, 4, 1, 6, 8, 2, 7, 5, 3, 0, 9, 8, 3, 7, 4, 0, 5, 6, 2, 1, 4, 9, 3, 8, 6, 2, 5, 7, 1, 0, 8, 9, 0, 7, 4, 1, 6, 5, 3, 2}
                };

                List<int[]> fgReelSet = new List<int[]>()
                {
                    new int[] {4, 5, 3, 7, 9, 1, 8, 6, 0, 7, 8, 1, 6, 4, 0, 9, 5, 1, 6, 8, 3, 9, 5, 1, 6, 4, 2, 8, 7, 3, 4, 5, 2, 6, 7, 3, 8, 9, 0, 7, 4, 3, 9, 5, 1, 6, 8, 2, 4, 5, 0, 6, 7, 3, 5, 9, 1, 8, 9, 2, 6, 4, 0, 7, 8, 2, 5, 4, 1, 9, 7, 3, 6, 9, 2, 4, 5, 3, 6, 8, 0, 9, 7, 2, 8, 6, 0, 5, 4, 3, 8, 7, 1, 9, 8, 2, 7, 9, 0, 6, 7, 3, 4, 5, 1, 9, 8, 2, 1, 4, 7, 0, 5, 8, 2, 9, 6, 3, 5, 9, 3, 7, 6, 1, 4, 8, 0, 9, 4, 1, 6, 8, 2, 7, 5, 0, 9, 8, 3, 7, 4, 0, 5, 6, 1, 2, 4, 9, 3, 8, 6, 2, 5, 7, 0, 8, 9, 0, 7, 4, 1, 6, 5, 2, 4, 7, 0, 5, 8, 2, 9, 6, 1, 3, 5, 9, 3, 7, 6, 1, 4, 8, 2, 9, 4, 1, 6, 8, 2, 7, 5, 3, 0, 9, 8, 3, 7, 4, 0, 5, 6, 2, 3, 4, 9, 3, 8, 6, 2, 5, 7, 1, 8, 9, 0, 7, 4, 1, 6, 5, 3, 2},
                    new int[] {4, 5, 3, 7, 9, 1, 8, 6, 0, 7, 8, 1, 6, 4, 0, 9, 5, 1, 6, 8, 3, 9, 5, 1, 6, 4, 2, 8, 7, 3, 4, 5, 2, 6, 7, 3, 8, 9, 0, 7, 4, 3, 9, 5, 1, 6, 8, 2, 4, 5, 0, 6, 7, 3, 5, 9, 1, 8, 9, 2, 6, 4, 0, 7, 8, 2, 5, 4, 1, 9, 7, 3, 6, 9, 2, 4, 5, 3, 6, 8, 0, 9, 7, 2, 8, 6, 0, 5, 4, 3, 8, 7, 1, 9, 8, 2, 7, 9, 0, 6, 7, 3, 4, 5, 1, 9, 8, 2, 1, 4, 7, 0, 5, 8, 2, 9, 6, 3, 5, 9, 3, 7, 6, 1, 4, 8, 0, 9, 4, 1, 6, 8, 2, 7, 5, 0, 9, 8, 3, 7, 4, 0, 5, 6, 1, 2, 4, 9, 3, 8, 6, 2, 5, 7, 0, 8, 9, 0, 7, 4, 1, 6, 5, 2, 4, 7, 0, 5, 8, 2, 9, 6, 1, 3, 5, 9, 3, 7, 6, 1, 4, 8, 2, 9, 4, 1, 6, 8, 2, 7, 5, 3, 0, 9, 8, 3, 7, 4, 0, 5, 6, 2, 3, 4, 9, 3, 8, 6, 2, 5, 7, 1, 8, 9, 0, 7, 4, 1, 6, 5, 3, 2},
                    new int[] {4, 5, 3, 7, 9, 1, 8, 6, 0, 7, 8, 1, 6, 4, 0, 9, 5, 1, 6, 8, 3, 9, 5, 1, 6, 4, 2, 8, 7, 3, 4, 5, 2, 6, 7, 3, 8, 9, 0, 7, 4, 3, 9, 5, 1, 6, 8, 2, 4, 5, 0, 6, 7, 3, 5, 9, 1, 8, 9, 2, 6, 4, 0, 7, 8, 2, 5, 4, 1, 9, 7, 3, 6, 9, 2, 4, 5, 3, 6, 8, 0, 9, 7, 2, 8, 6, 0, 5, 4, 3, 8, 7, 1, 9, 8, 2, 7, 9, 0, 6, 7, 3, 4, 5, 1, 9, 8, 2, 1, 4, 7, 0, 5, 8, 2, 9, 6, 3, 5, 9, 3, 7, 6, 1, 4, 8, 0, 9, 4, 1, 6, 8, 2, 7, 5, 0, 9, 8, 3, 7, 4, 0, 5, 6, 1, 2, 4, 9, 3, 8, 6, 2, 5, 7, 0, 8, 9, 0, 7, 4, 1, 6, 5, 2, 4, 7, 0, 5, 8, 2, 9, 6, 1, 3, 5, 9, 3, 7, 6, 1, 4, 8, 2, 9, 4, 1, 6, 8, 2, 7, 5, 3, 0, 9, 8, 3, 7, 4, 0, 5, 6, 2, 3, 4, 9, 3, 8, 6, 2, 5, 7, 1, 8, 9, 0, 7, 4, 1, 6, 5, 3, 2},
                    new int[] {4, 5, 3, 7, 9, 1, 8, 6, 0, 7, 8, 1, 6, 4, 0, 9, 5, 1, 6, 8, 3, 9, 5, 1, 6, 4, 2, 8, 7, 3, 4, 5, 2, 6, 7, 3, 8, 9, 0, 7, 4, 3, 9, 5, 1, 6, 8, 2, 4, 5, 0, 6, 7, 3, 5, 9, 1, 8, 9, 2, 6, 4, 0, 7, 8, 2, 5, 4, 1, 9, 7, 3, 6, 9, 2, 4, 5, 3, 6, 8, 0, 9, 7, 2, 8, 6, 0, 5, 4, 3, 8, 7, 1, 9, 8, 2, 7, 9, 0, 6, 7, 3, 4, 5, 1, 9, 8, 2, 1, 4, 7, 0, 5, 8, 2, 9, 6, 3, 5, 9, 3, 7, 6, 1, 4, 8, 0, 9, 4, 1, 6, 8, 2, 7, 5, 0, 9, 8, 3, 7, 4, 0, 5, 6, 1, 2, 4, 9, 3, 8, 6, 2, 5, 7, 0, 8, 9, 0, 7, 4, 1, 6, 5, 2, 4, 7, 0, 5, 8, 2, 9, 6, 1, 3, 5, 9, 3, 7, 6, 1, 4, 8, 2, 9, 4, 1, 6, 8, 2, 7, 5, 3, 0, 9, 8, 3, 7, 4, 0, 5, 6, 2, 3, 4, 9, 3, 8, 6, 2, 5, 7, 1, 8, 9, 0, 7, 4, 1, 6, 5, 3, 2},
                    new int[] {4, 5, 3, 7, 9, 1, 8, 6, 0, 7, 8, 1, 6, 4, 0, 9, 5, 1, 6, 8, 3, 9, 5, 1, 6, 4, 2, 8, 7, 3, 4, 5, 2, 6, 7, 3, 8, 9, 0, 7, 4, 3, 9, 5, 1, 6, 8, 2, 4, 5, 0, 6, 7, 3, 5, 9, 1, 8, 9, 2, 6, 4, 0, 7, 8, 2, 5, 4, 1, 9, 7, 3, 6, 9, 2, 4, 5, 3, 6, 8, 0, 9, 7, 2, 8, 6, 0, 5, 4, 3, 8, 7, 1, 9, 8, 2, 7, 9, 0, 6, 7, 3, 4, 5, 1, 9, 8, 2, 1, 4, 7, 0, 5, 8, 2, 9, 6, 3, 5, 9, 3, 7, 6, 1, 4, 8, 0, 9, 4, 1, 6, 8, 2, 7, 5, 0, 9, 8, 3, 7, 4, 0, 5, 6, 1, 2, 4, 9, 3, 8, 6, 2, 5, 7, 0, 8, 9, 0, 7, 4, 1, 6, 5, 2, 4, 7, 0, 5, 8, 2, 9, 6, 1, 3, 5, 9, 3, 7, 6, 1, 4, 8, 2, 9, 4, 1, 6, 8, 2, 7, 5, 3, 0, 9, 8, 3, 7, 4, 0, 5, 6, 2, 3, 4, 9, 3, 8, 6, 2, 5, 7, 1, 8, 9, 0, 7, 4, 1, 6, 5, 3, 2}
                };

                List<int[]> payTable = new List<int[]>()
                {
                    new int[] {500,200,30,10,0,0},
                    new int[] {300,100,20,5,0,0},
                    new int[] {200,60,15,0,0,0},
                    new int[] {150,50,10,0,0,0},

                    new int[] {100,40,8,0,0,0},
                    new int[] {100,40,8,0,0,0},
                    new int[] {80,30,6,0,0,0},
                    new int[] {80,30,6,0,0,0},
                    new int[] {60,20,5,0,0,0},
                    new int[] {60,20,5,0,0,0},
                    new int[] {1000,400,150,20,0,0},
                    new int[] {0,0,0,0,0,0}
                };

                int[] scatPay = new int[] { 0, 0, 10, 40, 200, 1000 };

                List<int[]> payLines = new List<int[]>()
                {
                    new int[] { 1, 1, 1, 1, 1 },
                    new int[] { 2, 2, 2, 2, 2 },
                    new int[] { 0, 0, 0, 0, 0 },
                    new int[] { 0, 1, 2, 1, 0 },
                    new int[] { 2, 1, 0, 1, 2 },

                    new int[] { 1, 0, 0, 0, 1 },
                    new int[] { 1, 2, 2, 2, 1 },
                    new int[] { 2, 1, 1, 1, 2 },
                    new int[] { 0, 1, 1, 1, 0 },
                    new int[] { 0, 0, 1, 2, 2 },

                    new int[] { 2, 2, 1, 0, 0 },
                    new int[] { 1, 2, 1, 2, 1 },
                    new int[] { 1, 0, 1, 0, 1 },
                    new int[] { 2, 1, 2, 1, 2 },
                    new int[] { 0, 1, 0, 1, 0 }
                };

                //base game
                int[] weiNWild = new int[] { 2100, 820, 200, 32, 6, 1 };
                int[] weiPosWild = new int[] { 1, 1, 4, 5, 6 };

                int[] weiNScat = new int[] { 24000, 4800, 900, 180, 11, 1 };
                int[] weiPosScat = new int[] { 6, 5, 4, 3, 2 };

                //free games
                int[] weiMulti = new int[] { 2, 10, 50, 150 };
                int[] weiLevelScat = new int[] { 0, 3, 7, 12 };
                
                List<int[]> weiNWildFG = new List<int[]>()
                {
                    new[] { 2400, 3800, 490, 50, 9, 1 },
                    new[] { 2700, 3600, 430, 45, 8, 1 },
                    new[] { 3100, 3400, 380, 40, 7, 1 },
                    new[] { 3900, 3200, 340, 35, 6, 1 }
                };

                int[] weiPosWildFG = new int[] { 7, 7, 8, 9, 9 };

                List<int[]> weiNScatFG = new List<int[]>()
                {
                    new[] { 11000, 16000, 3000, 180, 10, 1 },
                    new[] { 70000, 20000, 2800, 560, 10, 1 },
                    new[] { 80000, 23000, 2600, 620, 10, 1 },
                    new[] { 90000, 25000, 2400, 680, 10, 1 }
                };

                int[] weiPosScatFG = new int[] { 1, 1, 1, 1, 1 };
                

                //buy bonus
                /*
                weiNScat = new int[] { 0, 0, 0, 180, 11, 1 };

                List<int[]> weiNWildFG = new List<int[]>()
                {
                    new[] { 2000, 5000, 1500, 180, 14, 1 },
                    new[] { 2200, 4800, 1200, 140, 12, 1 },
                    new[] { 2500, 4600, 1000, 110, 10, 1 },
                    new[] { 2900, 4400, 900, 90, 8, 1 }
                };
                int[] weiPosWildFG = new int[] { 1, 1, 1, 1, 1 };

                List<int[]> weiNScatFG = new List<int[]>()
                {
                    new[] { 10000, 18000, 3000, 200, 10, 1 },
                    new[] { 60000, 22000, 2800, 560, 10, 1 },
                    new[] { 70000, 25000, 2600, 640, 10, 1 },
                    new[] { 80000, 27000, 2400, 710, 10, 1 }
                };
                int[] weiPosScatFG = new int[] { 1, 1, 1, 1, 1 };
                */

                int[] bgRanStop = new int[] { 0, 0, 0, 0, 0 };

                int[] bgScreen = new int[nx * ny];

                int[] bgHit = new int[payLines.Count + 1];
                int[] bgWin = new int[payLines.Count + 1];

                List<int[]> tmpreelset = new List<int[]>();
                List<int[]> tmpfgreelset = new List<int[]>();

                bool bonTriggered = false;

                List<int[]> fgScreens;
                List<int> fgLevel;
                List<int[]> fgRanStop;
                List<int[]> fgHits;
                List<int[]> fgWins;
                List<int> fgLen;
                List<int> fgScatters;


                int nrun = 0;
                while (nrun < ntot)
                {
                    nrun++;

                    int tot_win = 0;
                    
                    //creating screen
                    if (rnd.Next(2) == 0) tmpreelset = bgReelSet1;
                    else tmpreelset = bgReelSet2;

                    
                    for (int i = 0; i < nx; i++)
                    {
                        bgRanStop[i] = rnd.Next(tmpreelset[i].Length);

                        for (int j = 0; j < ny; j++)
                        {
                            int ind = (bgRanStop[i] + j) % tmpreelset[i].Length;
                            bgScreen[ny * i + j] = tmpreelset[i][ind];
                        }
                    }

                    //inserting wilds
                    int nWilds = GetWeiRand(weiNWild, rnd);
                    //nWilds = 0;
                    int[] twlwei = new int[nx];
                    for (int i = 0; i < nx; i++) twlwei[i] = weiPosWild[i];

                    for (int i = 0; i < nWilds; i++)
                    {
                        int col = GetWeiRand(twlwei, rnd);
                        int row = rnd.Next(ny);
                        int ind = ny * col + row;
                        bgScreen[ind] = 10;
                        twlwei[col] = 0;
                    }

                    //inserting scatters
                    int nScat = GetWeiRand(weiNScat, rnd);
                    //nScat = 5;
                    int[] tscwei = new int[nx];
                    for (int i = 0; i < nx; i++) tscwei[i] = weiPosScat[i];

                    for (int i = 0; i < nScat; i++)
                    {
                        int col = GetWeiRand(tscwei, rnd);
                        int row = rnd.Next(ny);
                        int ind = ny * col + row;
                        bgScreen[ind] = 11;
                        tscwei[col] = 0;
                    }

                    //getting wins
                    for (int i = 0; i < payLines.Count + 1; i++) { bgHit[i] = 0; bgWin[i] = 0; }

                    for (int i = 0; i < payLines.Count; i++)
                    {

                        //line combination
                        int[] combination = new int[nx];
                        for (int j = 0; j < nx; j++) combination[j] = bgScreen[ny * j + payLines[i][j]];

                        //symbol type
                        int symbol = 10;
                        int firstwilds = 0;
                        for (int j = 0; j < nx; j++)
                        {
                            if (combination[j] != 10)
                            {
                                symbol = combination[j];
                                break;
                            }
                            firstwilds++;
                        }

                        for (int j = 0; j < nx; j++)
                        {
                            if (combination[j] == symbol) bgHit[i]++;
                            else if (combination[j] == 10) bgHit[i]++;
                            else break;
                        }

                        if (firstwilds == 4)
                        {
                            if (payTable[symbol][5 - bgHit[i]] <= 400) { symbol = 10; bgHit[i] = 4; }
                        }
                        else if (firstwilds == 3)
                        {
                            if (payTable[symbol][5 - bgHit[i]] <= 150) { symbol = 10; bgHit[i] = 3; }
                        }
                        else if (firstwilds == 2)
                        {
                            if (payTable[symbol][5 - bgHit[i]] <= 20) { symbol = 10; bgHit[i] = 2; }
                        }

                        int prize = payTable[symbol][5 - bgHit[i]];
                        if (prize > 0) bgWin[i] = prize;
                        else bgHit[i] = 0;

                        //ccc[symbol, bgHit[i]]++;

                    }

                    //scatter check
                    bonTriggered = false;

                    bgHit[payLines.Count] = nScat;
                    bgWin[payLines.Count] = scatPay[nScat];

                    int bg_win = bgWin.Sum();
                    tot_win += bg_win;
                    
                    if (nScat > 2) bonTriggered = true;
                    //bonTriggered = true;

                    //free games
                    int fg_win = 0;
                    fgScreens = new List<int[]>();
                    fgLevel = new List<int>();
                    fgRanStop = new List<int[]>();
                    fgHits = new List<int[]>();
                    fgWins = new List<int[]>();
                    fgLen = new List<int>();
                    fgScatters = new List<int>();

                    if (bonTriggered)
                    {
                        int nfs = 5;

                        int fgmulti;
                        int level = 0;
                        int nscat = 0;


                        int k = 0;
                        while (k < nfs)
                        {

                            fgLevel.Add(level);

                            fgmulti = weiMulti[level];

                            tmpfgreelset = fgReelSet;

                            //skrin and reelstops
                            int[] transtop = new int[nx];
                            int[] tskrin = new int[nx * ny];

                            for (int i = 0; i < nx; i++)
                            {
                                transtop[i] = rnd.Next(tmpfgreelset[i].Length);

                                for (int j = 0; j < ny; j++)
                                {
                                    int ind = (transtop[i] + j) % tmpfgreelset[i].Length;
                                    tskrin[ny * i + j] = tmpfgreelset[i][ind];
                                }

                            }

                            //inserting wilds
                            int nwilds = GetWeiRand(weiNWildFG[level], rnd);
                            //nwilds = 0;
                            int[] twlweifg = new int[nx];
                            for (int i = 0; i < nx; i++) twlweifg[i] = weiPosWildFG[i];

                            for (int i = 0; i < nwilds; i++)
                            {
                                int col = GetWeiRand(twlweifg, rnd);
                                int row = rnd.Next(ny);
                                int ind = ny * col + row;
                                tskrin[ind] = 10;
                                twlweifg[col] = 0;
                            }

                            //inserting scatters
                            int nnscat = GetWeiRand(weiNScatFG[level], rnd);
                            //nnscat = 5;
                            int[] tscweifg = new int[nx];
                            for (int i = 0; i < nx; i++) tscweifg[i] = weiPosScatFG[i];

                            for (int i = 0; i < nnscat; i++)
                            {
                                int col = GetWeiRand(tscweifg, rnd);
                                int row = rnd.Next(ny);
                                int ind = ny * col + row;
                                tskrin[ind] = 11;
                                tscweifg[col] = 0;
                            }

                            fgScreens.Add(tskrin);
                            fgRanStop.Add(transtop);

                            //fg wins
                            int[] thit = new int[payLines.Count + 1];
                            int[] twin = new int[payLines.Count + 1];

                            for (int i = 0; i < payLines.Count; i++)
                            {
                                //line combination
                                int[] combination = new int[nx];
                                for (int j = 0; j < nx; j++) combination[j] = tskrin[ny * j + payLines[i][j]];

                                //symbol type
                                int symbol = 10;
                                int firstwilds = 0;
                                for (int j = 0; j < nx; j++)
                                {
                                    if (combination[j] != 10)
                                    {
                                        symbol = combination[j];
                                        break;
                                    }
                                    firstwilds++;
                                }

                                for (int j = 0; j < nx; j++)
                                {
                                    if (combination[j] == symbol) thit[i]++;
                                    else if (combination[j] == 10) thit[i]++;
                                    else break;
                                }

                                if (firstwilds == 4)
                                {
                                    if (payTable[symbol][5 - thit[i]] <= 400) { symbol = 10; thit[i] = 4; }
                                }
                                else if (firstwilds == 3)
                                {
                                    if (payTable[symbol][5 - thit[i]] <= 150) { symbol = 10; thit[i] = 3; }
                                }
                                else if (firstwilds == 2)
                                {
                                    if (payTable[symbol][5 - thit[i]] <= 20) { symbol = 10; thit[i] = 2; }
                                }

                                int prize = payTable[symbol][5 - thit[i]];
                                if (prize > 0) twin[i] = prize;
                                else thit[i] = 0;


                                bool wldhit = false;
                                for (int w = 0; w < thit[i]; w++) if (combination[w] == 10) { wldhit = true; break; }

                                if (wldhit) twin[i] = fgmulti * prize;

                            }


                            //free games length
                            fgLen.Add(nfs);


                            //scatter hits
                            //int nnscat = 0;
                            if (nnscat > 2) nfs += 5;

                            thit[payLines.Count] = nnscat;
                            twin[payLines.Count] = scatPay[nnscat];


                            //level up check
                            nscat += nnscat;

                            if (level < 2)
                            {
                                if (nscat >= weiLevelScat[level + 1])
                                {
                                    level++;
                                    nfs += 5;
                                    if (nscat >= weiLevelScat[level + 1]) { level++; nfs += 5; }
                                }
                            }
                            else if (level == 2)
                            {
                                if (nscat >= weiLevelScat[3]) { level++; nfs += 5; }
                            }


                            fgHits.Add(thit);
                            fgWins.Add(twin);
                            fgScatters.Add(nscat);

                            fg_win += twin.Sum();

                            //max win check
                            tot_win += twin.Sum();
                            if (tot_win > 10000 * bet)
                            {
                                tot_win = 10000 * bet;
                                fg_win = tot_win - bg_win;
                                cap++;
                                break;
                            }
                            


                            //ith free spin
                            k++;
                        }

                    }


                    total += tot_win;
                    bg_total += bg_win;
                    fg_total += fg_win;

                    //std_total += (bg_win + fg_win) * (bg_win + fg_win);
                    
                    if (nrun % (ntot / 10) == 0) Console.WriteLine(nrun);
                }

                result[5, ip] = (float)total / ntot / bet;
                result[0, ip] = (float)bg_total / ntot / bet;
                result[1, ip] = (float)fg_total / ntot / bet;
                result[2, ip] = (float)Math.Sqrt( (float)std_total / ntot / bet / bet );

                result[3, ip] = (float)plus / ntot / bet;
                result[4, ip] = cap;

            });


            float[] res_end = new float[10];

            for (int i = 0; i < 10; i++)
            {
                res_end[i] = 0;
                for (int j = 0; j < 6; j++) res_end[i] += result[i, j];
            }

            Debug.Print("rtp = " + (res_end[5]) / 6);
            Debug.Print("rtp bg = " + res_end[0] / 6);
            Debug.Print("rtp fg = " + res_end[1] / 6);
            Debug.Print("stdev = " + res_end[2] / 6);
            Debug.Print("plus = " + res_end[3] / 6);
            Debug.Print("ncap = " + res_end[4]);

        }

        static int GetWeiRand(int[] weights, Random rnd)
        {
            if (weights == null || weights.Length == 0) return -1;

            int t = 0;
            int i;
            for (i = 0; i < weights.Length; i++) t += weights[i];

            if (t == 0) return -1;

            //rnd = new Random(DateTime.Now.Millisecond);
            float r = rnd.Next(t); ;

            int s = 0;
            for (i = 0; i < weights.Length; i++)
            {
                s += weights[i];
                if (s > r) return i;
            }

            return -1;
        }

    }
}

