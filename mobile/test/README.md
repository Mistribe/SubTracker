# Testing Approach

This document outlines the testing approach for the Recurrent Payment Tracker application.

## Unit Tests

### Models

The `Payment` model is tested in `test/models/payment_test.dart`. These tests verify:
- Total amount spent calculation for monthly payments
- Total amount spent calculation for annual payments
- Total amount spent calculation for new payments
- Formatted total amount spent

### Providers

The `PaymentProvider` is tested in `test/providers/payment_provider_test.dart`. These tests verify:
- Active payments count
- Total amount spent calculation

## Mock Repository

For testing the `PaymentProvider`, we use a mock repository that implements the `PaymentRepository` interface. This allows us to test the provider without relying on actual persistence.

Key points about the mock repository:
- It maintains an in-memory list of payments
- It returns a copy of the payments list to avoid unintended modifications
- It implements all methods of the `PaymentRepository` interface

## Hive Initialization

Since our models use Hive for persistence, we need to initialize Hive for testing. This is done in the `setUpAll` method of each test file:

```
setUpAll(() async {
  // Set up a temporary directory for Hive
  final tempDir = await Directory.systemTemp.createTemp('hive_test');
  Hive.init(tempDir.path);

  // Register adapters
  if (!Hive.isAdapterRegistered(0)) {
    Hive.registerAdapter(PaymentAdapter());
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

Since the `PaymentProvider` methods are async, we need to make our test methods async as well. We also need to wait for the provider to initialize before running the tests:

```
setUp(() async {
  mockRepository = MockPaymentRepository();
  provider = PaymentProvider(paymentRepository: mockRepository);

  // Wait for the provider to initialize
  await Future.delayed(Duration.zero);
});
```

## Widget Testing

A placeholder widget test is included in `test/widget_test.dart`. This can be expanded in the future to test the UI components of the application.
