import 'package:uuid/uuid.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/foundation.dart';
import '../models/family_member.dart';
import '../models/subscription.dart';
import '../models/subscription_payment.dart';
import '../models/label.dart';
import '../models/currency.dart';
import '../models/subscription_state.dart';
import '../repositories/subscription_repository.dart';
import '../repositories/settings_repository.dart';
import '../repositories/label_repository.dart';
import '../services/currency_converter.dart';
import '../providers/sync_provider.dart';

var uuid = Uuid();

// Enum for subscription sorting options
enum SubscriptionSortOption {
  none,
  nameAsc,
  nameDesc,
  nextPaymentAsc,
  nextPaymentDesc,
}

enum SubscriptionFilterOption {
  labels,
  showInactive,
  hideInactive,
  familyMembers,
  payer,
}

// Special value for family common account in payer filter
const String kFamilyCommonAccountId = 'family_common_account';

class SubscriptionProvider with ChangeNotifier {
  final SubscriptionRepository subscriptionRepository;
  final SettingsRepository? settingsRepository;
  final LabelRepository labelRepository;
  final List<String> _selectedLabelIds = [];
  List<Subscription> _subscriptions = [];
  bool _showInactiveSubscriptions = false;
  SubscriptionSortOption _sortOption = SubscriptionSortOption.none;
  String? _selectedFamilyMemberId;
  String? _selectedPayerFamilyMemberId;

  SubscriptionProvider({
    required this.subscriptionRepository,
    required this.labelRepository,
    this.settingsRepository,
  }) {
    // Load subscriptions from repository
    _loadPayments();
    // Load settings if repository is provided
    if (settingsRepository != null) {
      _loadSettings();
    }
  }

  // Load settings from repository
  void _loadSettings() {
    final settings = settingsRepository!.getSettings();
    _defaultCurrency = settings.defaultCurrency;
  }

  // Load subscriptions from repository
  Future<void> _loadPayments() async {
    _subscriptions = subscriptionRepository.getAll();
    notifyListeners();
  }

  // Getter for the subscriptions list
  List<Subscription> get subscriptions => List.unmodifiable(_subscriptions);

  // Getter for filtered subscriptions with optimized implementation
  List<Subscription> get filteredSubscriptions {
    // Apply all filters in a single pass
    final filtered = _subscriptions.where((subscription) {
      // Filter by active status
      if (!_showInactiveSubscriptions &&
          subscription.state == SubscriptionState.ended) {
        return false;
      }

      // Filter by selected labels
      if (_selectedLabelIds.isNotEmpty &&
          !subscription.labelIds.any(
            (labelId) => _selectedLabelIds.contains(labelId),
          )) {
        return false;
      }

      // Filter by selected family member
      if (_selectedFamilyMemberId != null &&
          !subscription.userFamilyMemberIds.contains(_selectedFamilyMemberId)) {
        return false;
      }

      // Filter by selected payer
      if (_selectedPayerFamilyMemberId != null) {
        // Special case for "Family (common account)"
        if (_selectedPayerFamilyMemberId == kFamilyCommonAccountId) {
          // Show subscriptions with no specific payer
          if (subscription.payerFamilyMemberId != null) {
            return false;
          }
        } else if (subscription.payerFamilyMemberId == null ||
            subscription.payerFamilyMemberId != _selectedPayerFamilyMemberId) {
          return false;
        }
      }

      return true;
    }).toList();

    // Apply sorting
    if (_sortOption != SubscriptionSortOption.none) {
      _applySorting(filtered);
    }

    return List.unmodifiable(filtered);
  }

  // Helper method to apply sorting
  void _applySorting(List<Subscription> subscriptions) {
    switch (_sortOption) {
      case SubscriptionSortOption.nameAsc:
        subscriptions.sort(
          (a, b) => a.name.toLowerCase().compareTo(b.name.toLowerCase()),
        );
        break;
      case SubscriptionSortOption.nameDesc:
        subscriptions.sort(
          (a, b) => b.name.toLowerCase().compareTo(a.name.toLowerCase()),
        );
        break;
      case SubscriptionSortOption.nextPaymentAsc:
        subscriptions.sort((a, b) {
          if (a.state != SubscriptionState.active &&
              b.state != SubscriptionState.active) {
            return 0;
          }
          if (a.state != SubscriptionState.active) return 1;
          if (b.state != SubscriptionState.active) return -1;
          return a.nextPaymentDate.compareTo(b.nextPaymentDate);
        });
        break;
      case SubscriptionSortOption.nextPaymentDesc:
        subscriptions.sort((a, b) {
          if (a.state != SubscriptionState.active &&
              b.state != SubscriptionState.active) {
            return 0;
          }
          if (a.state != SubscriptionState.active) return 1;
          if (b.state != SubscriptionState.active) return -1;
          return b.nextPaymentDate.compareTo(a.nextPaymentDate);
        });
        break;
      default:
        break;
    }
  }

  // Getter for the show inactive subscriptions setting
  bool get showInactiveSubscriptions => _showInactiveSubscriptions;

  // Setter for the show inactive subscriptions setting
  set showInactiveSubscriptions(bool value) {
    _showInactiveSubscriptions = value;
    notifyListeners();
  }

  // Toggle the show inactive subscriptions setting
  void toggleShowInactiveSubscriptions() {
    _showInactiveSubscriptions = !_showInactiveSubscriptions;
    notifyListeners();
  }

  // Getter for the current sort option
  SubscriptionSortOption get sortOption => _sortOption;

  // Setter for the sort option
  set sortOption(SubscriptionSortOption option) {
    _sortOption = option;
    notifyListeners();
  }

  // Getter for all labels
  List<Label> get labels => List.unmodifiable(labelRepository.getAll());

  // Getter for default labels
  List<Label> get defaultLabels => labelRepository.getDefaultLabels();

  // Getter for custom labels
  List<Label> get customLabels => labelRepository.getCustomLabels();

  // Getter for selected label IDs
  List<String> get selectedLabelIds => List.unmodifiable(_selectedLabelIds);

  // Getter for selected labels (Label objects)
  List<Label> get selectedLabels {
    // Get the latest labels from the repository
    final allLabels = labelRepository.getAll();
    return allLabels
        .where((label) => _selectedLabelIds.contains(label.id))
        .toList();
  }

  // Toggle a label in the filter (add if not present, remove if present)
  void toggleLabelFilter(String labelId) {
    final isSelected = _selectedLabelIds.contains(labelId);

    if (isSelected) {
      _selectedLabelIds.remove(labelId);
    } else {
      _selectedLabelIds.add(labelId);
    }

    notifyListeners();
  }

  // Clear all label filters
  void clearLabelFilters() {
    if (_selectedLabelIds.isNotEmpty) {
      _selectedLabelIds.clear();
      notifyListeners();
    }
  }

  // Getter for selected family member ID
  String? get selectedFamilyMemberId => _selectedFamilyMemberId;

  // Setter for selected family member ID
  set selectedFamilyMemberId(String? familyMemberId) {
    _selectedFamilyMemberId = familyMemberId;
    notifyListeners();
  }

  // Clear family member filter
  void clearFamilyMemberFilter() {
    if (_selectedFamilyMemberId != null) {
      _selectedFamilyMemberId = null;
      notifyListeners();
    }
  }

  // Getter for selected payer family member ID
  String? get selectedPayerFamilyMemberId => _selectedPayerFamilyMemberId;

  // Setter for selected payer family member ID
  set selectedPayerFamilyMemberId(String? payerFamilyMemberId) {
    _selectedPayerFamilyMemberId = payerFamilyMemberId;
    notifyListeners();
  }

  // Clear payer family member filter
  void clearPayerFamilyMemberFilter() {
    if (_selectedPayerFamilyMemberId != null) {
      _selectedPayerFamilyMemberId = null;
      notifyListeners();
    }
  }

  // Clear all filters
  void clearAllFilters() {
    clearLabelFilters();
    clearFamilyMemberFilter();
    clearPayerFamilyMemberFilter();
  }

  // Calculate total monthly cost in the default currency (USD)
  double get totalMonthlyCost {
    return _subscriptions.fold(
      0,
      (sum, subscription) => sum + subscription.monthlyCost,
    );
  }

  // Calculate total monthly cost in the selected currency
  double get totalMonthlyCostInSelectedCurrency {
    return _subscriptions.fold(0.0, (sum, subscription) {
      final detail = subscription.getLastPaymentDetail();
      if (detail.state != SubscriptionState.active) return sum;

      // Convert the monthly cost from the subscription's currency to the selected currency
      double convertedCost = CurrencyConverter.convert(
        subscription.monthlyCost,
        detail.currency,
        _defaultCurrency,
      );

      return sum + convertedCost;
    });
  }

  // Calculate total annual cost in the selected currency
  double get totalAnnualCostInSelectedCurrency {
    return totalMonthlyCostInSelectedCurrency * 12;
  }

  // Format the monthly cost with the selected currency symbol
  String get formattedMonthlyCost {
    return CurrencyConverter.formatAmountWithCurrency(
      totalMonthlyCostInSelectedCurrency,
      _defaultCurrency,
    );
  }

  // Format the annual cost with the selected currency symbol
  String get formattedAnnualCost {
    return CurrencyConverter.formatAmountWithCurrency(
      totalAnnualCostInSelectedCurrency,
      _defaultCurrency,
    );
  }

  // Get the count of active subscriptions (excluding stopped subscriptions)
  int get activePaymentsCount {
    return _subscriptions
        .where((subscription) => subscription.state == SubscriptionState.active)
        .length;
  }

  // Get the count of active subscriptions that have not started yet
  int get notStartedPaymentsCount {
    return _subscriptions
        .where(
          (subscription) => subscription.state == SubscriptionState.notStarted,
        )
        .length;
  }

  // Calculate total amount spent across all subscriptions
  double get totalAmountSpent {
    return _subscriptions.fold(
      0,
      (sum, subscription) => sum + subscription.totalAmountSpent,
    );
  }

  String _defaultCurrency = Currency.USD.code;

  // Getter for default currency
  String get defaultCurrency => _defaultCurrency;

  // Setter for default currency
  set defaultCurrency(String currency) {
    _defaultCurrency = currency;

    // Persist settings if repository is available
    if (settingsRepository != null) {
      settingsRepository!.updateDefaultCurrency(currency);
    }

    notifyListeners();
  }

  // Get the default currency as a Currency enum
  Currency get defaultCurrencyEnum => Currency.fromCode(_defaultCurrency);

  // Set the default currency using a Currency enum
  set defaultCurrencyEnum(Currency currency) {
    defaultCurrency = currency.code;
  }

  // Add a new subscription
  Future<void> addPayment(
    String name,
    double price,
    int months,
    DateTime startDate, {
    DateTime? endDate,
    String? currency,
    List<String>? labelIds,
    List<String>? userFamilyMemberIds,
    String? payerFamilyMemberId,
  }) async {
    // Validate that a kid is not set as a payer
    // Note: We can't validate if a kid is set as a payer here anymore since we only have the ID
    // This validation should be done before calling this method

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
      labelIds: labelIds ?? [],
      userFamilyMemberIds: userFamilyMemberIds,
      payerFamilyMemberId: payerFamilyMemberId,
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
  Future<void> updateSubscription(
    String id,
    String name,
    List<String> labelIds,
    List<String> familyMemberIds,
    String? payerId,
  ) async {
    final index = _subscriptions.indexWhere(
      (subscription) => subscription.id == id,
    );

    if (index >= 0) {
      final subscription = _subscriptions[index];

      _subscriptions[index] = subscription.copyWith(
        name: name,
        userFamilyMemberIds: familyMemberIds,
        payerFamilyMemberId: payerId,
        labelIds: labelIds,
        updatedAt: DateTime.now(),
      );

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

      // Create new subscription payment
      final newPayment = SubscriptionPayment(
        id: _generateId(),
        price: newPrice,
        startDate: effectiveDate,
        endDate: endDate,
        months: months ?? previousDetail.months,
        currency: currency ?? previousDetail.currency,
      );

      // Add new subscription detail
      subscription.addPaymentDetail(newPayment);

      // Persist to storage
      await subscriptionRepository.createPayment(subscription.id, newPayment);

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
      final detailCurrency =
          currency ?? (existingDetail?.currency ?? _defaultCurrency);

      // Create updated payment
      final updatedPayment = SubscriptionPayment(
        id: subscriptionDetailId,
        price: newPrice,
        startDate: startDate,
        endDate: endDate,
        months: months,
        currency: detailCurrency,
      );

      // Update the subscription
      subscription.setPaymentDetail(updatedPayment);

      // Persist to storage
      await subscriptionRepository.updatePayment(
        subscriptionId,
        updatedPayment,
      );

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

      // Get the current payment detail before updating it
      final currentDetail = subscription.getLastPaymentDetail();

      // Update the subscription
      subscription.endCurrentPaymentDetail(effectiveStopDate);

      // Get the updated payment detail
      final updatedDetail = subscription.getLastPaymentDetail();

      // Persist to storage
      await subscriptionRepository.updatePayment(
        subscription.id,
        updatedDetail,
      );

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

      // Create new subscription payment
      final newPayment = SubscriptionPayment(
        id: _generateId(),
        price: price ?? previousDetail.price,
        startDate: reactivationDate,
        endDate: endDate,
        months: months ?? previousDetail.months,
        currency: currency ?? previousDetail.currency,
      );

      // Add to subscription
      subscription.addPaymentDetail(newPayment);

      // Persist to storage
      await subscriptionRepository.createPayment(subscription.id, newPayment);

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

      await subscriptionRepository.deletePayment(subscription.id, paymentId);

      notifyListeners();
    }
  }

  // Generate a unique ID for a new subscription
  String _generateId() {
    return uuid.v7();
  }

  // Label management methods have been moved to LabelProvider

  // Update subscription labels
  Future<void> updateSubscriptionLabels(
    String subscriptionId,
    List<String> labelIds,
  ) async {
    final index = _subscriptions.indexWhere((s) => s.id == subscriptionId);
    if (index >= 0) {
      final subscription = _subscriptions[index];
      final updatedSubscription = subscription.copyWith(labelIds: labelIds);

      _subscriptions[index] = updatedSubscription;
      await subscriptionRepository.update(updatedSubscription);

      notifyListeners();
    }
  }

  // Update subscription family members
  Future<void> updateSubscriptionFamilyMembers(
    String subscriptionId, {
    List<String>? userFamilyMemberIds,
    String? payerFamilyMemberId,
  }) async {
    // Validate that a kid is not set as a payer
    // Note: We can't validate if a kid is set as a payer here anymore since we only have the ID
    // This validation should be done before calling this method

    final index = _subscriptions.indexWhere((s) => s.id == subscriptionId);
    if (index >= 0) {
      final subscription = _subscriptions[index];
      final updatedSubscription = subscription.copyWith(
        userFamilyMemberIds: userFamilyMemberIds,
        payerFamilyMemberId: payerFamilyMemberId,
      );

      _subscriptions[index] = updatedSubscription;
      await subscriptionRepository.update(updatedSubscription);

      notifyListeners();
    }
  }
}
