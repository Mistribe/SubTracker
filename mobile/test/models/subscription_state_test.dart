import 'package:flutter_test/flutter_test.dart';
import 'package:subscription_tracker/models/subscription_state.dart';

void main() {
  group('SubscriptionState', () {
    test('enum has correct values', () {
      expect(SubscriptionState.values.length, equals(3));
      expect(SubscriptionState.values, contains(SubscriptionState.active));
      expect(SubscriptionState.values, contains(SubscriptionState.notStarted));
      expect(SubscriptionState.values, contains(SubscriptionState.ended));
    });

    test('enum values have correct indices', () {
      expect(SubscriptionState.active.index, equals(0));
      expect(SubscriptionState.notStarted.index, equals(1));
      expect(SubscriptionState.ended.index, equals(2));
    });
  });
}