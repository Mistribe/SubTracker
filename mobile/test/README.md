# Testing Approach

This document outlines the testing approach for the Subscription Tracker application.

## Unit Tests

### Models

The `Subscription` model is tested in `test/models/subscription_test.dart`. These tests verify:
- Total amount spent calculation for monthly subscriptions
- Total amount spent calculation for annual subscriptions
- Total amount spent calculation for new subscriptions
- Formatted total amount spent

### Providers

The `SubscriptionProvider` is tested in `test/providers/subscription_provider_test.dart`. These tests verify:
- Active subscriptions count
- Total amount spent calculation

## Mock Repository

For testing the `SubscriptionProvider`, we use a mock repository that implements the `SubscriptionRepository` interface. This allows us to test the provider without relying on actual persistence.

Key points about the mock repository:
- It maintains an in-memory list of subscriptions
- It returns a copy of the subscriptions list to avoid unintended modifications
- It implements all methods of the `SubscriptionRepository` interface

## Hive Initialization

Since our models use Hive for persistence, we need to initialize Hive for testing. This is done in the `setUpAll` method of each test file:

```
setUpAll(() async {
  // Set up a temporary directory for Hive
  final tempDir = await Directory.systemTemp.createTemp('hive_test');
  Hive.init(tempDir.path);

  // Register adapters
  if (!Hive.isAdapterRegistered(0)) {
    Hive.registerAdapter(SubscriptionAdapter());
  }
  if (!Hive.isAdapterRegistered(1)) {
    Hive.registerAdapter(PriceChangeAdapter());
  }
});
```

And we clean up in the `tearDownAll` method:

```
tearDownAll(() async {
  // Clean up after all tests
  await Hive.close();
});
```

## Async Testing

Since the `SubscriptionProvider` methods are async, we need to make our test methods async as well. We also need to wait for the provider to initialize before running the tests:

```
setUp(() async {
  mockRepository = MockSubscriptionRepository();
  provider = SubscriptionProvider(subscriptionRepository: mockRepository);

  // Wait for the provider to initialize
  await Future.delayed(Duration.zero);
});
```

## Widget Testing

A placeholder widget test is included in `test/widget_test.dart`. This can be expanded in the future to test the UI components of the application.
