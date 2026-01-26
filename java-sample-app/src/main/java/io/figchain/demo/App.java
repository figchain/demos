package io.figchain.demo;

import java.nio.file.Path;
import java.util.Optional;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;

import io.figchain.client.FigChainClient;
import io.figchain.client.FigChainClientBuilder;
import io.figchain.client.transport.FcAuthenticationException;
import io.figchain.client.transport.FcAuthorizationException;
import io.figchain.schema.test.test;

public class App {

    private static final String TEST_KEY = "test";

    public static void main(String[] args) throws Exception {
        boolean once = false;
        for (String arg : args) {
            if ("--once".equals(arg)) {
                once = true;
                break;
            }
        }

        // Build the client
        FigChainClientBuilder builder = new FigChainClientBuilder()
                .fromConfig(Path.of("client-config.json"));
        FigChainClient client = builder.build();

        try {
            // Register a listener for the test fig
            client.registerListener(TEST_KEY, test.class, g -> {
                System.out.println("Received updated test: " + g.getTest());
            });

            // Start the client and fetch the test fig
            System.out.println("Starting client and fetching initial data...");
            client.start().get(10, TimeUnit.SECONDS);
            System.out.println("Client started");
            Optional<test> test = client.getFig(TEST_KEY, test.class);
            test.ifPresent(t -> System.out.println("Fetched test: " + t.getTest()));
            if (!test.isPresent()) {
                System.out.println("Test not found");
            }

            if (once) {
                System.out.println("Exiting because --once was specified.");
                client.stop();
                return;
            }

            // Wait for the user to press Enter to exit
            System.out.println("Press Enter to exit...");
            System.in.read();
            client.stop();
        } catch (ExecutionException e) {
            // Handle authentication and authorization errors
            if (e.getCause() instanceof FcAuthenticationException || e.getCause() instanceof FcAuthorizationException) {
                System.err.println("Authentication failed. Please check authPrivateKey and ensure that you have access to the requested workspace.");
            } else {
                System.err.println("Execution failed: " + e.getMessage());
            }
        }
    }
}
