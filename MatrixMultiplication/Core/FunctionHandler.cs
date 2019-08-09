using System;
using System.Collections.Generic;
using System.IO;
using System.Net.Http;
using System.Runtime.Serialization.Json;
using System.Text;
using System.Threading.Tasks;
using MatrixMul.Core.Interfaces;
using MatrixMul.Core.Model;
using Newtonsoft.Json;

namespace MatrixMul.Core
{
    public class FunctionHandler
    {
        private readonly IMatrixMulRepository datastore;
        private readonly Random rnd = new Random();

        private ComputationHandler cHandler = new ComputationHandler();

        public FunctionHandler(IMatrixMulRepository datastore)
        {
            this.datastore = datastore;
        }

        public string CreateMatrix(int size, int maxValue = 50000)
        {
            var uid = Guid.NewGuid().ToString();
            var mtxA = Util.GenerateMatrix(size, (x, y) => rnd.Next(maxValue));
            var mtxB = Util.GenerateMatrix(size, (x, y) => rnd.Next(maxValue));

            var c = new MatrixCalculation
            {
                A = mtxA,
                B = mtxB
            };

            datastore.StoreCalculation(uid, c);
            return uid;
        }

        public void SerialMultiply(string id)
        {
            var calc = datastore.GetCalculation(id);

            var result = cHandler.SerialMultiply(calc);
            datastore.StoreResultMatrix(id, result);
        }

        public void ScheduleMultiplicationTasks(string id, int workerCount)
        {
            var calc = datastore.GetCalculation(id);

            foreach (var workerTasks in cHandler.BuildTasks(workerCount, calc))
            {
                Console.WriteLine($"Scheduling worker {workerTasks.Key}");
                this.datastore.StoreComputationTasksForWorker(id, workerTasks.Key, workerTasks.Value.ToArray());
            }
        }

        public void ParallelMultiplyWorker(string id, int workerId)
        {
            var tasks = datastore.GetComputationTasksForWorker(id, workerId);
            var calc = datastore.GetCalculation(id);
            var results = cHandler.PerformCalculations(workerId, new List<ComputationTask>(tasks), calc);

            datastore.StoreComputationResults(id, workerId, results.ToArray());
        }

        public void BuildResultMatrix(string id, int workerCount)
        {
            var calc = datastore.GetCalculation(id);

            List<List<ComputationResult>> results = new List<List<ComputationResult>>();
            for (int i = 0; i < workerCount; i++)
            {
                var data = datastore.GetComputationResults(id, i);
                results.Add(new List<ComputationResult>(data));
            }

            var rMatrix = cHandler.BuildResultMatrix(calc, results);
            datastore.StoreResultMatrix(id, rMatrix);
        }

        public Report GenerateReport(string callbackUrl, long start, string id, int workerCount)
        {
            var calc = datastore.GetCalculation(id);
            var result = datastore.GetResultMatrix(id);

            var deltask = Task.Run(() =>
            {
                // Cleanup
                datastore.DeleteCalculation(id);
                datastore.DeleteResultMatrix(id);

                for (int i = 0; i < workerCount; i++)
                {
                    datastore.DeleteComputationResults(id, i);
                    datastore.DeleteComputationTasks(id, i);
                }
            });

            var doneTs = Util.GetUnixTimestamp();

            var report = new Report
            {
                Size = result.Size,
                EndTimestamp = doneTs,
                StartTimestamp = start,
                ResultMatrix = result.ToMatrixInfo(),
                InputMatrixA = calc.A.ToMatrixInfo(),
                InputMatrixB = calc.B.ToMatrixInfo()
            };

            if (callbackUrl != null)
            {
                var client = new HttpClient();
                var data = JsonConvert.SerializeObject(report);
                Console.WriteLine(data);
                Task.WaitAll(client.PostAsync(callbackUrl,
                    new StringContent(data, Encoding.Default, "application/json")));
            }

            Task.WaitAll(deltask);

            return report;
        }
    }
}