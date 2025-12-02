package io.figchain.demo;

import java.util.Optional;
import java.util.Set;
import java.util.UUID;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.TimeUnit;

import io.figchain.client.FigChainClient;
import io.figchain.client.FigChainClientBuilder;
import io.figchain.client.transport.FcAuthenticationException;
import io.figchain.client.transport.FcAuthorizationException;
import io.figchain.schema.greeting.Greeting;

public class App {

    private static final String GREETING_KEY = "my-greeting";

    public static void main(String[] args) throws Exception {
        // Fetch client params from the environment
        if (System.getenv("FIGCHAIN_CREDENTIAL") == null || System.getenv("FIGCHAIN_ENVIRONMENT_ID") == null) {
            System.err.println("Environment variables FIGCHAIN_CREDENTIAL and FIGCHAIN_ENVIRONMENT_ID must be set");
            System.exit(64);
        }

        // Build the client
        FigChainClientBuilder builder = new FigChainClientBuilder();
        builder.withClientSecret(System.getenv("FIGCHAIN_CREDENTIAL"));
        builder.withEnvironmentId(UUID.fromString(System.getenv("FIGCHAIN_ENVIRONMENT_ID")));
        builder.withNamespaces(Set.of("default"));
        FigChainClient client = builder.build();

        try {
            // Register a listener for the greeting fig
            client.registerListener(GREETING_KEY, Greeting.class, g -> {
                System.out.println("Received updated greeting: " + g.getMessage());
            });

            // Start the client and fetch the greeting fig
            System.out.println("Starting client and fetching initial data...");
            client.start().get(10, TimeUnit.SECONDS);
            System.out.println("Client started");
            Optional<Greeting> greeting = client.getFig(GREETING_KEY, Greeting.class);
            greeting.ifPresent(g -> System.out.println("Fetched greeting: " + g.getMessage()));
            if (!greeting.isPresent()) {
                System.out.println("Greeting not found");
            }

            // Wait for the user to press Enter to exit
            System.out.println("Press Enter to exit...");
            System.in.read();
            client.stop();
        } catch (ExecutionException e) {
            // Handle authentication and authorization errors
            if (e.getCause() instanceof FcAuthenticationException || e.getCause() instanceof FcAuthorizationException) {
                System.err.println("Authentication failed: " + e.getCause().getMessage());
            } else {
                System.err.println("Execution failed: " + e.getMessage());
            }
        }
    }
}
