﻿using System;
using System.Runtime.Serialization;
using MatrixMul.Core;
using Newtonsoft.Json.Linq;

namespace MatrixMul.IBMCloud
{
    public class CreateMatrix
    {
        public JObject Main(JObject args)
        {
            try
            {
                var size = args.ContainsKey("size") ? int.Parse(args["size"].ToString()) : 50;
                var max = args.ContainsKey("max") ? int.Parse(args["max"].ToString()) : 5000;
                args["hasCallback"] = args.ContainsKey("callback");
                args["startTime"] = Util.GetUnixTimestamp();

                var repo = new S3Repository(args);
                var hndlr = new FunctionHandler(repo);

                var id = hndlr.CreateMatrix(size, max);

                args["id"] = id;
                args["size"] = size;
                args["max"] = max;

                Console.WriteLine(args.ToString());

                return args;
            }
            catch (Exception e)
            {
                var j = new JObject();
                j["error"] = e.ToString();
                return j;
            }
        }
    }
}