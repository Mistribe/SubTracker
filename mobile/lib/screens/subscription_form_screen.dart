import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import 'package:subscription_tracker/models/subscription_state.dart';
import '../providers/subscription_provider.dart';
import '../providers/family_provider.dart';
import '../providers/label_provider.dart';
import '../models/subscription.dart';
import '../models/label.dart';
import '../models/currency.dart';
import '../models/family_member.dart';
// Widget imports
import '../widgets/subscription_basic_info_section.dart';
import '../widgets/subscription_details_section.dart';
import '../widgets/subscription_family_section.dart';
import '../widgets/subscription_labels_section.dart';
import '../widgets/color_selection_dialog.dart';

class SubscriptionFormScreen extends StatefulWidget {
  final Subscription? subscription;

  const SubscriptionFormScreen({super.key, this.subscription});

  @override
  State<SubscriptionFormScreen> createState() => _SubscriptionFormScreenState();
}

class _SubscriptionFormScreenState extends State<SubscriptionFormScreen> {
  final _formKey = GlobalKey<FormState>();
  late final TextEditingController _nameController;
  late final TextEditingController _priceController;
  late final TextEditingController _customMonthsController;
  late final TextEditingController _recurrenceCountController;
  late final TextEditingController _freeTrialMonthsController;
  late int _months;
  String _selectedDuration = '1 month'; // Default value
  late DateTime _startDate;
  DateTime? _endDate;
  late String _selectedCurrency;
  late List<Label> _selectedLabels;
  bool _isLabelsExpanded = false; // Track if labels section is expanded
  bool _isFamilyExpanded = false; // Track if family section is expanded
  bool _useRecurrenceCount = false;
  bool _useFreeTrial = false; // Track if free trial is enabled
  int _freeTrialMonths = 1; // Default free trial months
  bool _isEditMode = false;
  bool _isActiveSubscription = true;
  List<FamilyMember> _userFamilyMembers = [];
  FamilyMember? _payerFamilyMember;
  bool _payedByJointAccount = false;
  bool _isFormModified = false; // Track if form has been modified

  // Get currencies from the Currency enum
  final List<String> _currencies = Currency.codes;

  @override
  void initState() {
    super.initState();
    _isEditMode = widget.subscription != null;

    if (_isEditMode) {
      // Edit mode - initialize with subscription data
      final subscription = widget.subscription!;
      _isActiveSubscription = subscription.state == SubscriptionState.active;

      _nameController = TextEditingController(text: subscription.name);

      _payedByJointAccount = subscription.payedByJointAccount;
      // Get the current subscription detail
      final currentDetail = subscription.getLastPaymentDetail();
      _priceController = TextEditingController(
        text: currentDetail.price.toString(),
      );

      _months = currentDetail.months;
      _startDate = currentDetail.startDate;
      _endDate = currentDetail.endDate;
      if (_endDate != null) {
        int recurrenceCount = _calculateRecurrenceCount(
          currentDetail.startDate,
          currentDetail.endDate!,
          currentDetail.months,
        );
        _recurrenceCountController = TextEditingController(
          text: recurrenceCount.toString(),
        );
        _useRecurrenceCount = true;
      } else {
        _recurrenceCountController = TextEditingController(text: '1');
      }
      _selectedCurrency = currentDetail.currency;

      // Initialize free trial data
      _freeTrialMonths = currentDetail.freeTrialMonths;
      _useFreeTrial = _freeTrialMonths > 0;
      _freeTrialMonthsController = TextEditingController(
        text: _freeTrialMonths > 0 ? _freeTrialMonths.toString() : '1',
      );

      // Fetch labels from provider
      WidgetsBinding.instance.addPostFrameCallback((_) {
        final labelProvider = Provider.of<LabelProvider>(
          context,
          listen: false,
        );
        setState(() {
          _selectedLabels = subscription.labelIds
              .map((id) => labelProvider.getLabel(id) ?? Label.empty())
              .where((label) => label.id.isNotEmpty)
              .toList();
        });
      });

      // Fetch family members from provider
      WidgetsBinding.instance.addPostFrameCallback((_) {
        final familyProvider = Provider.of<FamilyProvider>(
          context,
          listen: false,
        );
        setState(() {
          _userFamilyMembers = subscription.familyMemberIds
              .map(
                (id) =>
                    familyProvider.getFamilyMemberById(id) ??
                    FamilyMember.empty(),
              )
              .where((member) => member.id.isNotEmpty)
              .toList();

          _payerFamilyMember = subscription.payerId != null
              ? familyProvider.getFamilyMemberById(subscription.payerId)
              : null;
        });
      });

      // Set the selected duration based on months
      if (_months == 1) {
        _selectedDuration = '1 month';
      } else if (_months == 3) {
        _selectedDuration = '3 months';
      } else if (_months == 6) {
        _selectedDuration = '6 months';
      } else if (_months == 12) {
        _selectedDuration = '12 months';
      } else {
        _selectedDuration = 'Custom';
      }

      _customMonthsController = TextEditingController(
        text: _selectedDuration == 'Custom' ? _months.toString() : '',
      );
    } else {
      // Add mode - initialize with default values
      _nameController = TextEditingController();
      _priceController = TextEditingController();
      _customMonthsController = TextEditingController();
      _recurrenceCountController = TextEditingController(text: '1');
      _freeTrialMonthsController = TextEditingController(text: '1');
      _months = 1;
      _startDate = DateTime.now();
      _selectedLabels = [];

      // Initialize with the default currency from the provider
      WidgetsBinding.instance.addPostFrameCallback((_) {
        final provider = Provider.of<SubscriptionProvider>(
          context,
          listen: false,
        );
        setState(() {
          _selectedCurrency = provider.defaultCurrency;
        });
      });

      // Fetch family members from provider
      WidgetsBinding.instance.addPostFrameCallback((_) {
        final familyProvider = Provider.of<FamilyProvider>(
          context,
          listen: false,
        );
        setState(() {
          final selectedFamily = familyProvider.selectedFamily;
          if (selectedFamily != null) {
            _payedByJointAccount = selectedFamily.haveJointAccount;
          } else {
            _payedByJointAccount = false;
          }
        });
      });

      // Default currency until provider loads
      _selectedCurrency = Currency.USD.code;
    }

    // Add listeners to controllers to detect changes
    _nameController.addListener(_onFormChanged);
    _priceController.addListener(_onFormChanged);
    _customMonthsController.addListener(_onFormChanged);
    _recurrenceCountController.addListener(_onFormChanged);
  }

  // Called when any form field changes
  void _onFormChanged() {
    if (!_isFormModified) {
      setState(() {
        _isFormModified = true;
      });
    }
  }

  @override
  void dispose() {
    _nameController.dispose();
    _priceController.dispose();
    _customMonthsController.dispose();
    _recurrenceCountController.dispose();
    super.dispose();
  }

  // Date picker for start date
  Future<void> _selectStartDate(BuildContext context) async {
    // Don't allow editing start date for inactive subscriptions in edit mode
    if (_isEditMode && !_isActiveSubscription) return;

    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: _startDate,
      firstDate: DateTime(2000),
      lastDate: DateTime(2101),
    );
    if (picked != null && picked != _startDate) {
      setState(() {
        _startDate = picked;
        _isFormModified = true;
      });
    }
  }

  // Handle duration selection
  void _handleDurationChange(String? value) {
    // Don't allow editing recurrence for inactive subscriptions in edit mode
    if (_isEditMode && !_isActiveSubscription) return;

    if (value == null) return;

    setState(() {
      _selectedDuration = value;
      _isFormModified = true;

      // Update months based on selection
      if (value == '1 month') {
        _months = 1;
      } else if (value == '3 months') {
        _months = 3;
      } else if (value == '6 months') {
        _months = 6;
      } else if (value == '12 months') {
        _months = 12;
      }
      // For custom, we'll use the value from the text field
    });
  }

  // Get the appropriate label for the duration
  String _getDurationLabel() {
    switch (_selectedDuration) {
      case '1 month':
        return 'months';
      case '3 months':
        return 'quarters';
      case '6 months':
        return 'half-years';
      case '12 months':
        return 'years';
      case 'Custom':
        final customMonths =
            int.tryParse(_customMonthsController.text) ?? _months;
        if (customMonths == 12) {
          return 'years';
        } else if (customMonths == 6) {
          return 'half-years';
        } else if (customMonths == 3) {
          return 'quarters';
        } else {
          return 'periods of $customMonths months';
        }
      default:
        return 'recurrences';
    }
  }

  // Handle custom months change
  void _handleCustomMonthsChange(String value) {
    // Don't allow editing recurrence for inactive subscriptions in edit mode
    if (_isEditMode && !_isActiveSubscription) return;

    if (value.isEmpty) return;

    final parsedValue = int.tryParse(value);
    if (parsedValue != null && parsedValue > 0) {
      setState(() {
        _months = parsedValue;
        _isFormModified = true;
      });
    }
  }

  int _calculateRecurrenceCount(
    DateTime startDate,
    DateTime endDate,
    int months,
  ) {
    int yearDifference = endDate.year - startDate.year;
    int monthDifference = endDate.month - startDate.month;

    return ((monthDifference + (yearDifference * 12)) / months).truncate();
  }

  // Calculate end date based on start date, months, and recurrence count
  DateTime _calculateEndDate() {
    final recurrenceCount = int.tryParse(_recurrenceCountController.text) ?? 1;
    final year = (_months / 12).truncate();
    final month = _months % 12;

    return DateTime(
      _startDate.year + (year * recurrenceCount),
      _startDate.month + (month * recurrenceCount),
      _startDate.day,
    );
  }

  Future<void> _submitForm() async {
    if (_formKey.currentState!.validate()) {
      // Store context in local variables before async operations
      final currentContext = context;
      final messenger = ScaffoldMessenger.of(currentContext);
      final navigator = Navigator.of(currentContext);

      final name = _nameController.text.trim();
      final price = double.tryParse(_priceController.text) ?? 0.0;

      // For custom duration, get months from text field
      if (_selectedDuration == 'Custom') {
        final customMonths = int.tryParse(_customMonthsController.text);
        if (customMonths != null && customMonths > 0) {
          _months = customMonths;
        }
      }

      // Calculate the end date based on recurrence count if toggle is on (only in add mode)
      if (!_isEditMode && _useRecurrenceCount) {
        final recurrenceCount =
            int.tryParse(_recurrenceCountController.text) ?? 1;
        if (recurrenceCount > 0) {
          _endDate = _calculateEndDate();
        }
      }
      try {
        // Show loading indicator
        messenger.showSnackBar(
          SnackBar(
            content: Text(
              _isEditMode
                  ? 'Updating subscription...'
                  : 'Adding subscription...',
            ),
            duration: const Duration(seconds: 1),
          ),
        );

        final subscriptionProvider = Provider.of<SubscriptionProvider>(
          context,
          listen: false,
        );

        final familyProvider = Provider.of<FamilyProvider>(
          context,
          listen: false,
        );

        if (_isEditMode) {
          // Update mode
          final subscription = widget.subscription!;

          // Update the subscription name and labels
          await subscriptionProvider.updateSubscription(
            subscription.id,
            name,
            _selectedLabels.map((label) => label.id).toList(),
            _userFamilyMembers.map((member) => member.id).toList(),
            _payerFamilyMember?.id,
            payedByJointAccount: _payedByJointAccount,
            familyId:
                _payedByJointAccount == true ||
                    _userFamilyMembers.isNotEmpty ||
                    _payerFamilyMember != null
                ? familyProvider.selectedFamilyId
                : null,
          );

          // Update the subscription detail if it's active
          if (_isActiveSubscription) {
            final currentDetail = subscription.getLastPaymentDetail();
            await subscriptionProvider.updateSubscriptionPayment(
              subscription.id,
              currentDetail.id,
              price,
              _months,
              _startDate,
              endDate: _endDate,
              currency: _selectedCurrency,
              freeTrialMonths: _useFreeTrial ? _freeTrialMonths : 0,
            );
          }
        } else {
          // Add mode
          await subscriptionProvider.addPayment(
            name,
            price,
            _months,
            _startDate,
            endDate: _endDate,
            currency: _selectedCurrency,
            labelIds: _selectedLabels.map((label) => label.id).toList(),
            userFamilyMemberIds: _userFamilyMembers
                .map((member) => member.id)
                .toList(),
            payerFamilyMemberId: _payerFamilyMember?.id,
            payedByJointAccount: _payedByJointAccount,
            familyId:
                _payedByJointAccount == true ||
                    _userFamilyMembers.isNotEmpty ||
                    _payerFamilyMember != null
                ? familyProvider.selectedFamilyId
                : null,
            freeTrialMonths: _useFreeTrial ? _freeTrialMonths : 0,
          );
        }

        // Show success message
        messenger.showSnackBar(
          SnackBar(
            content: Text(
              _isEditMode
                  ? 'Subscription updated successfully'
                  : 'Subscription added successfully',
            ),
            duration: const Duration(seconds: 2),
          ),
        );

        // Reset form modified flag
        _isFormModified = false;

        // Close the form
        navigator.pop();
      } catch (e) {
        // Show error message
        messenger.showSnackBar(
          SnackBar(
            content: Text(
              'Error ${_isEditMode ? 'updating' : 'adding'} subscription: ${e.toString()}',
            ),
            backgroundColor: Theme.of(context).colorScheme.error,
            duration: const Duration(seconds: 3),
          ),
        );
      }
    }
  }


  // Show dialog for adding a custom label
  void _showAddLabelDialog() {
    final nameController = TextEditingController();
    final selectedColor = ValueNotifier<Color>(Theme.of(context).colorScheme.primary);

    showDialog(
      context: context,
      builder: (context) => ColorSelectionDialog(
        nameController: nameController,
        selectedColor: selectedColor,
        onAddLabel: (name, colorHex) {
          Provider.of<LabelProvider>(
            context,
            listen: false,
          ).addLabel(name, colorHex);
        },
      ),
    );
  }

  // Show confirmation dialog when trying to go back with unsaved changes
  Future<void> _showUnsavedChangesDialog() async {
    final navigator = Navigator.of(context);
    if (!_isFormModified) {
      navigator.pop(); // No changes, just go back
      return;
    }

    // Show confirmation dialog
    final result = await showDialog<String>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Unsaved Changes'),
        content: const Text(
          'You have unsaved changes. What would you like to do?',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop('cancel'),
            // Don't discard changes
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () {
              Navigator.of(context).pop('discard'); // Discard changes
            },
            child: const Text('Discard modification'),
          ),
          TextButton(
            onPressed: () {
              Navigator.of(context).pop('save'); // Save changes
            },
            child: const Text('Save'),
          ),
        ],
      ),
    );

    if (result == 'discard') {
      navigator.pop(); // Go back without saving
    } else if (result == 'save') {
      _submitForm(); // Save changes
    }
    // If result is 'cancel' or null, do nothing and stay on the form
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(
          _isEditMode ? 'Edit Subscription' : 'Add Subscription',
        ),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: _showUnsavedChangesDialog,
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.check),
            onPressed: _submitForm,
            tooltip: _isEditMode ? 'Update Subscription' : 'Add Subscription',
          ),
        ],
      ),
      body: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Form(
            key: _formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [

                // Basic Information Section
                SubscriptionBasicInfoSection(
                  nameController: _nameController,
                  priceController: _priceController,
                  selectedCurrency: _selectedCurrency,
                  currencies: _currencies,
                  isEditMode: _isEditMode,
                  isActiveSubscription: _isActiveSubscription,
                  onCurrencyChanged: (String? newValue) {
                    if (newValue != null) {
                      setState(() {
                        _selectedCurrency = newValue;
                        _isFormModified = true;
                      });
                    }
                  },
                ),

                if (!_isEditMode || (_isEditMode && _isActiveSubscription)) ...[
                  const SizedBox(height: 24),

                  // Subscription Details Section
                  SubscriptionDetailsSection(
                    selectedDuration: _selectedDuration,
                    customMonthsController: _customMonthsController,
                    recurrenceCountController: _recurrenceCountController,
                    freeTrialMonthsController: _freeTrialMonthsController,
                    startDate: _startDate,
                    useRecurrenceCount: _useRecurrenceCount,
                    useFreeTrial: _useFreeTrial,
                    isEditMode: _isEditMode,
                    isActiveSubscription: _isActiveSubscription,
                    onDurationChanged: _handleDurationChange,
                    onCustomMonthsChanged: _handleCustomMonthsChange,
                    onSelectStartDate: _selectStartDate,
                    onRecurrenceCountToggled: (value) {
                      setState(() {
                        _useRecurrenceCount = value;
                        _isFormModified = true;
                        if (!value) {
                          _recurrenceCountController.text = '1';
                        }
                      });
                    },
                    onFreeTrialToggled: (value) {
                      setState(() {
                        _useFreeTrial = value;
                        _isFormModified = true;
                        if (!value) {
                          _freeTrialMonths = 0;
                        } else {
                          _freeTrialMonths = int.tryParse(_freeTrialMonthsController.text) ?? 1;
                        }
                      });
                    },
                    onFreeTrialMonthsChanged: (value) {
                      setState(() {
                        _freeTrialMonths = int.tryParse(value) ?? 1;
                        _isFormModified = true;
                      });
                    },
                    durationLabel: _getDurationLabel(),
                  ),
                ],

                const SizedBox(height: 24),

                // Family Section
                Consumer<FamilyProvider>(
                  builder: (context, familyMemberProvider, _) {
                    return SubscriptionFamilySection(
                      userFamilyMembers: _userFamilyMembers,
                      allFamilyMembers: familyMemberProvider.familyMembers,
                      payerFamilyMember: _payerFamilyMember,
                      payedByJointAccount: _payedByJointAccount,
                      isFamilyExpanded: _isFamilyExpanded,
                      hasFamilyMembers: familyMemberProvider.hasFamilyMembers,
                      onFamilyExpandedToggle: (value) {
                        setState(() {
                          _isFamilyExpanded = value;
                        });
                      },
                      onFamilyMemberToggled: (member, selected) {
                        setState(() {
                          if (selected) {
                            _userFamilyMembers.add(member);
                          } else {
                            _userFamilyMembers.removeWhere(
                              (selectedMember) => selectedMember.id == member.id,
                            );
                          }
                          _isFormModified = true;
                        });
                      },
                      onJointAccountToggled: (value) {
                        setState(() {
                          _payedByJointAccount = value;
                          if (value) {
                            // If joint account is selected, clear the payer
                            _payerFamilyMember = null;
                          }
                          _isFormModified = true;
                        });
                      },
                      onPayerChanged: (String? value) {
                        setState(() {
                          _payerFamilyMember = value == null
                              ? null
                              : familyMemberProvider.getFamilyMemberById(value);
                          _isFormModified = true;
                        });
                      },
                    );
                  },
                ),

                const SizedBox(height: 24),

                // Labels Section
                Consumer<LabelProvider>(
                  builder: (context, labelProvider, _) {
                    return SubscriptionLabelsSection(
                      selectedLabels: _selectedLabels,
                      allLabels: labelProvider.labels,
                      isLabelsExpanded: _isLabelsExpanded,
                      onLabelsExpandedToggle: (value) {
                        setState(() {
                          _isLabelsExpanded = value;
                        });
                      },
                      onLabelToggled: (label, selected) {
                        setState(() {
                          if (selected) {
                            _selectedLabels.add(label);
                          } else {
                            _selectedLabels.removeWhere(
                              (selectedLabel) => selectedLabel.id == label.id,
                            );
                          }
                          _isFormModified = true;
                        });
                      },
                      onAddCustomLabel: _showAddLabelDialog,
                    );
                  },
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
