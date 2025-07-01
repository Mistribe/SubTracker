import 'package:uuid/uuid.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/foundation.dart';
import '../models/subscription.dart';
import '../models/subscription_payment.dart';
import '../repositories/subscription_repository.dart';
import '../services/currency_converter.dart';

var uuid = Uuid();

class SubscriptionProvider with ChangeNotifier {
  final SubscriptionRepository subscriptionRepository;
  List<Subscription> _subscriptions = [];

  SubscriptionProvider({required this.subscriptionRepository}) {
    // Load subscriptions from repository
    _loadPayments();
  }

  // Load subscriptions from repository
  Future<void> _loadPayments() async {
    _subscriptions = subscriptionRepository.getAll();
    notifyListeners();
  }

  // Getter for the subscriptions list
  List<Subscription> get subscriptions => List.unmodifiable(_subscriptions);

  // Calculate total monthly cost in the default currency (USD)
  double get totalMonthlyCost {
    return _subscriptions.fold(
      0,
      (sum, subscription) => sum + subscription.monthlyCost,
    );
  }

  // Calculate total monthly cost in the selected currency
  double get totalMonthlyCostInSelectedCurrency {
    return _subscriptions.fold(
      0.0,
      (sum, subscription) {
        final detail = subscription.getLastPaymentDetail();
        if (!detail.isActive) return sum;

        // Convert the monthly cost from the subscription's currency to the selected currency
        double convertedCost = CurrencyConverter.convert(
          subscription.monthlyCost,
          detail.currency,
          _defaultCurrency
        );

        return sum + convertedCost;
      },
    );
  }

  // Calculate total annual cost in the selected currency
  double get totalAnnualCostInSelectedCurrency {
    return totalMonthlyCostInSelectedCurrency * 12;
  }

  // Format the monthly cost with the selected currency symbol
  String get formattedMonthlyCost {
    return CurrencyConverter.formatAmountWithCurrency(
      totalMonthlyCostInSelectedCurrency,
      _defaultCurrency
    );
  }

  // Format the annual cost with the selected currency symbol
  String get formattedAnnualCost {
    return CurrencyConverter.formatAmountWithCurrency(
      totalAnnualCostInSelectedCurrency,
      _defaultCurrency
    );
  }

  // Get the count of active subscriptions (excluding stopped subscriptions)
  int get activePaymentsCount {
    return _subscriptions.where((subscription) => subscription.isActive).length;
  }

  // Calculate total amount spent across all subscriptions
  double get totalAmountSpent {
    return _subscriptions.fold(
      0,
      (sum, subscription) => sum + subscription.totalAmountSpent,
    );
  }

  String _defaultCurrency = 'USD';

  // Getter for default currency
  String get defaultCurrency => _defaultCurrency;

  // Setter for default currency
  set defaultCurrency(String currency) {
    _defaultCurrency = currency;
    notifyListeners();
  }

  // Add a new subscription
  Future<void> addPayment(
    String name,
    double price,
    int months,
    DateTime startDate, {
    DateTime? endDate,
    String? currency,
  }) async {
    // Create initial subscription history entry
    final initialSubscriptionPayment = [
      SubscriptionPayment(
        id: _generateId(),
        price: price,
        startDate: startDate,
        endDate: endDate,
        // Far future date
        months: months,
        currency: currency ?? _defaultCurrency,
      ),
    ];

    final subscription = Subscription(
      id: _generateId(),
      name: name,
      subscriptionPayments: initialSubscriptionPayment,
    );

    // Add to local list
    _subscriptions.add(subscription);

    // Persist to storage
    await subscriptionRepository.add(subscription);

    notifyListeners();
  }

  // Remove a subscription
  Future<void> removePayment(String id) async {
    _subscriptions.removeWhere((subscription) => subscription.id == id);

    // Remove from storage
    await subscriptionRepository.delete(id);

    notifyListeners();
  }

  // Update an existing subscription
  Future<void> updateSubscription(String id, String name) async {
    final index = _subscriptions.indexWhere(
      (subscription) => subscription.id == id,
    );

    if (index >= 0) {
      final subscription = _subscriptions[index];

      _subscriptions[index] = subscription.copyWith(name: name);

      // Persist to storage
      await subscriptionRepository.update(subscription);

      notifyListeners();
    }
  }

  // Add a price change at a specific date
  Future<void> addPaymentDetailEntry(
    String subscriptionId,
    double newPrice,
    DateTime effectiveDate, {
    DateTime? endDate,
    int? months,
    String? currency,
  }) async {
    final index = _subscriptions.indexWhere(
      (subscription) => subscription.id == subscriptionId,
    );

    if (index >= 0) {
      final subscription = _subscriptions[index];

      final previousDetail = subscription.getLastPaymentDetail();
      // Update the current subscription detail with end date
      subscription.setEndDateToCurrentPaymentDetail(effectiveDate);
      // Add new subscription detail
      subscription.addPaymentDetail(
        SubscriptionPayment(
          id: _generateId(),
          price: newPrice,
          startDate: effectiveDate,
          endDate: endDate,
          months: months ?? previousDetail.months,
          currency: currency ?? previousDetail.currency,
        ),
      );

      // Persist to storage
      await subscriptionRepository.update(subscription);

      notifyListeners();
    }
  }

  // Update a subscription history entry
  Future<void> updateSubscriptionPayment(
    String subscriptionId,
    String subscriptionDetailId,
    double newPrice,
    int months,
    DateTime startDate, {
    DateTime? endDate,
    String? currency,
  }) async {
    final index = _subscriptions.indexWhere(
      (subscription) => subscription.id == subscriptionId,
    );

    if (index >= 0) {
      final subscription = _subscriptions[index];

      // Get the existing payment detail to preserve the currency if not provided
      final existingDetail = subscription.findDetailById(subscriptionDetailId);
      final detailCurrency = currency ?? (existingDetail?.currency ?? _defaultCurrency);

      subscription.setPaymentDetail(
        SubscriptionPayment(
          id: subscriptionDetailId,
          price: newPrice,
          startDate: startDate,
          endDate: endDate,
          months: months,
          currency: detailCurrency,
        ),
      );

      // Persist to storage
      await subscriptionRepository.update(subscription);

      notifyListeners();
    }
  }

  // Stop a subscription
  Future<void> cancelCurrentSubscription(
    String subscriptionId, {
    DateTime? stopDate,
  }) async {
    final index = _subscriptions.indexWhere(
      (subscription) => subscription.id == subscriptionId,
    );

    if (index >= 0) {
      final subscription = _subscriptions[index];

      // If stopDate is not provided, use the last subscription date
      final effectiveStopDate = stopDate ?? subscription.lastPaymentDate;
      subscription.endCurrentPaymentDetail(effectiveStopDate);

      // Persist to storage
      await subscriptionRepository.update(subscription);

      notifyListeners();
    }
  }

  // Reactivate a subscription at a specific date
  Future<void> reactivatePayment(
    String subscriptionId,
    DateTime reactivationDate, {
    double? price,
    int? months,
    DateTime? endDate,
    String? currency,
  }) async {
    final index = _subscriptions.indexWhere(
      (subscription) => subscription.id == subscriptionId,
    );

    if (index >= 0) {
      final subscription = _subscriptions[index];

      final previousDetail = subscription.getLastPaymentDetail();

      subscription.addPaymentDetail(
        SubscriptionPayment(
          id: _generateId(),
          price: price ?? previousDetail.price,
          startDate: reactivationDate,
          endDate: endDate,
          months: months ?? previousDetail.months,
          currency: currency ?? previousDetail.currency,
        ),
      );

      // Persist to storage
      await subscriptionRepository.update(subscription);

      notifyListeners();
    }
  }

  // Remove a specific subscription payment
  Future<void> removeSubscriptionPayment(
    String subscriptionId,
    String paymentId,
  ) async {
    final index = _subscriptions.indexWhere(
      (subscription) => subscription.id == subscriptionId,
    );

    if (index >= 0) {
      final subscription = _subscriptions[index];

      // Remove the payment from the subscription
      subscription.removePaymentDetail(paymentId);

      // Persist to storage
      await subscriptionRepository.update(subscription);

      notifyListeners();
    }
  }

  // Generate a unique ID for a new subscription
  String _generateId() {
    return uuid.v7();
  }
}
