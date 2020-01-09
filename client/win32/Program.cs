/* https://github.com/aspnet/AspNetCore.Docs/blob/master/aspnetcore/fundamentals/servers/kestrel.md
 */

using System;
using System.Security.Cryptography.X509Certificates;
using System.Text.Json;
using Microsoft.AspNetCore.Hosting;
using Microsoft.AspNetCore.Server.Kestrel.Https;
using Microsoft.Extensions.Hosting;
using Microsoft.AspNetCore.Builder;
using Microsoft.Extensions.Logging;

namespace WorkerService
{
    public class Program
    {
        public static void Main(string[] args)
        {
            CreateHostBuilder(args).Build().Run();
        }

        private static String SerializeCertificate(X509Certificate2 certificate)
        {
            return String.Format("{0}",
                                    JsonSerializer.Serialize(new
                                    {
                                        Thumbprint = certificate.Thumbprint,
                                        SubjectName = certificate.Subject,
                                        NotBefore = certificate.NotBefore,
                                        NotAfter = certificate.NotAfter
                                    },
                                    new JsonSerializerOptions
                                    {
                                        WriteIndented = true
                                    }));
        }

        public static IHostBuilder CreateHostBuilder(string[] args) =>

            Host.CreateDefaultBuilder(args)
            .ConfigureWebHostDefaults(webBuilder =>
        {
            X509Store store = new X509Store(StoreLocation.CurrentUser);
            store.Open(OpenFlags.ReadOnly);
            var certificates = X509Certificate2UI.SelectFromCollection(
                store.Certificates,
                "test-client",
                "Select Certificate for Client API",
                X509SelectionFlag.SingleSelection);
            store.Close();

            var serverCert = certificates[0];
            Console.WriteLine("Listening Certificate:" + SerializeCertificate(serverCert));
            
            webBuilder.UseStartup<Startup>();

            webBuilder.ConfigureKestrel(kestrelServerOptions =>
                {
                    kestrelServerOptions.ConfigureHttpsDefaults(opt =>
                    {
                        opt.ServerCertificate = serverCert;
                        opt.ClientCertificateMode = ClientCertificateMode.RequireCertificate;
                        opt.ClientCertificateValidation = (certificate, chain, errors) =>
                        {
                            Startup.log.LogInformation("Remote Certificate:" + SerializeCertificate(certificate));

                            bool remoteCertificateAccepted = certificate.Thumbprint == "9F8DF580E2F1388282679678F4D5788FF4C4258C";  // some test certificate

                            if (!remoteCertificateAccepted)
                            {
                                Startup.log.LogError("Remote certificate with Thumbprint={0} DOES NOT ACCEPTED", certificate.Thumbprint);
                            }
                            else
                            {
                                Startup.log.LogDebug("Remote certificate with Thumbprint={0} ACCEPTED", certificate.Thumbprint);
                            }

                            return remoteCertificateAccepted;
                        };
                    });
                });

        });
    }
}
