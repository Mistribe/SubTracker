import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../providers/subscription_provider.dart';
import '../providers/family_member_provider.dart';
import '../models/subscription.dart';
import '../models/label.dart';
import '../models/currency.dart';
import '../models/family_member.dart';

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
  late int _months;
  String _selectedDuration = '1 month'; // Default value
  late DateTime _startDate;
  DateTime? _endDate;
  late String _selectedCurrency;
  late List<Label> _selectedLabels;
  bool _isLabelsExpanded = false; // Track if labels section is expanded
  bool _isFamilyExpanded = false; // Track if family section is expanded
  bool _useRecurrenceCount = false;
  bool _isEditMode = false;
  bool _isActiveSubscription = true;
  List<FamilyMember> _userFamilyMembers = [];
  FamilyMember? _payerFamilyMember;

  // Get currencies from the Currency enum
  final List<String> _currencies = Currency.codes;

  @override
  void initState() {
    super.initState();
    _isEditMode = widget.subscription != null;

    if (_isEditMode) {
      // Edit mode - initialize with subscription data
      final subscription = widget.subscription!;
      _isActiveSubscription = subscription.isActive;

      _nameController = TextEditingController(text: subscription.name);

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
      _selectedLabels = List<Label>.from(subscription.labels);
      _userFamilyMembers = List<FamilyMember>.from(subscription.userFamilyMembers);
      _payerFamilyMember = subscription.payerFamilyMember;

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

      // Default currency until provider loads
      _selectedCurrency = Currency.USD.code;
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
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              _isEditMode
                  ? 'Updating subscription...'
                  : 'Adding subscription...',
            ),
            duration: const Duration(seconds: 1),
          ),
        );

        final provider = Provider.of<SubscriptionProvider>(
          context,
          listen: false,
        );

        if (_isEditMode) {
          // Update mode
          final subscription = widget.subscription!;

          // Update the subscription name and labels
          await provider.updateSubscription(subscription.id, name);

          // Update the subscription labels
          await provider.updateSubscriptionLabels(
            subscription.id,
            _selectedLabels,
          );

          // Update the subscription family members
          await provider.updateSubscriptionFamilyMembers(
            subscription.id,
            userFamilyMembers: _userFamilyMembers,
            payerFamilyMember: _payerFamilyMember,
          );

          // Update the subscription detail if it's active
          if (_isActiveSubscription) {
            final currentDetail = subscription.getLastPaymentDetail();
            await provider.updateSubscriptionPayment(
              subscription.id,
              currentDetail.id,
              price,
              _months,
              _startDate,
              endDate: _endDate,
              currency: _selectedCurrency,
            );
          }
        } else {
          // Add mode
          await provider.addPayment(
            name,
            price,
            _months,
            _startDate,
            endDate: _endDate,
            currency: _selectedCurrency,
            labels: _selectedLabels,
            userFamilyMembers: _userFamilyMembers,
            payerFamilyMember: _payerFamilyMember,
          );
        }

        // Show success message
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              _isEditMode
                  ? 'Subscription updated successfully'
                  : 'Subscription added successfully',
            ),
            duration: const Duration(seconds: 2),
          ),
        );

        // Close the form
        Navigator.of(context).pop();
      } catch (e) {
        // Show error message
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              'Error ${_isEditMode ? 'updating' : 'adding'} subscription: ${e.toString()}',
            ),
            backgroundColor: Colors.red,
            duration: const Duration(seconds: 3),
          ),
        );
      }
    }
  }

  // Build label chips for selection
  List<Widget> _buildLabelChips() {
    final provider = Provider.of<SubscriptionProvider>(context);
    final allLabels = provider.labels;

    return allLabels.map((label) {
      final isSelected = _selectedLabels.any(
        (selectedLabel) => selectedLabel.id == label.id,
      );
      return FilterChip(
        label: Text(label.name),
        selected: isSelected,
        onSelected: (selected) {
          setState(() {
            if (selected) {
              _selectedLabels.add(label);
            } else {
              _selectedLabels.removeWhere(
                (selectedLabel) => selectedLabel.id == label.id,
              );
            }
          });
        },
        backgroundColor: Color(
          int.parse(label.color.substring(1, 7), radix: 16) + 0xFF000000,
        ).withOpacity(0.2),
        selectedColor: Color(
          int.parse(label.color.substring(1, 7), radix: 16) + 0xFF000000,
        ).withOpacity(0.7),
        checkmarkColor: Colors.white,
      );
    }).toList();
  }

  // Build family member chips for selection
  List<Widget> _buildFamilyMemberChips() {
    final familyMemberProvider = Provider.of<FamilyMemberProvider>(context);
    final allFamilyMembers = familyMemberProvider.familyMembers;

    return allFamilyMembers.map((member) {
      final isSelected = _userFamilyMembers.any(
        (selectedMember) => selectedMember.id == member.id,
      );
      return FilterChip(
        label: Text(member.name),
        selected: isSelected,
        onSelected: (selected) {
          setState(() {
            if (selected) {
              _userFamilyMembers.add(member);
            } else {
              _userFamilyMembers.removeWhere(
                (selectedMember) => selectedMember.id == member.id,
              );
            }
          });
        },
        backgroundColor: Theme.of(context).chipTheme.backgroundColor,
        selectedColor: Theme.of(context).colorScheme.primary.withOpacity(0.7),
        checkmarkColor: Colors.white,
      );
    }).toList();
  }

  // Show dialog for adding a custom label
  void _showAddLabelDialog() {
    final nameController = TextEditingController();
    final selectedColor = ValueNotifier<Color>(Colors.blue);

    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Add Custom Label'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            TextField(
              controller: nameController,
              decoration: const InputDecoration(
                labelText: 'Label Name',
                hintText: 'Enter a name for your label',
              ),
            ),
            const SizedBox(height: 16),
            const Text('Select Color:'),
            const SizedBox(height: 8),
            Wrap(
              spacing: 8,
              children: [
                _buildColorOption(selectedColor, Colors.red),
                _buildColorOption(selectedColor, Colors.pink),
                _buildColorOption(selectedColor, Colors.purple),
                _buildColorOption(selectedColor, Colors.deepPurple),
                _buildColorOption(selectedColor, Colors.indigo),
                _buildColorOption(selectedColor, Colors.blue),
                _buildColorOption(selectedColor, Colors.lightBlue),
                _buildColorOption(selectedColor, Colors.cyan),
                _buildColorOption(selectedColor, Colors.teal),
                _buildColorOption(selectedColor, Colors.green),
                _buildColorOption(selectedColor, Colors.lightGreen),
                _buildColorOption(selectedColor, Colors.lime),
                _buildColorOption(selectedColor, Colors.yellow),
                _buildColorOption(selectedColor, Colors.amber),
                _buildColorOption(selectedColor, Colors.orange),
                _buildColorOption(selectedColor, Colors.deepOrange),
                _buildColorOption(selectedColor, Colors.brown),
                _buildColorOption(selectedColor, Colors.grey),
                _buildColorOption(selectedColor, Colors.blueGrey),
              ],
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () {
              final name = nameController.text.trim();
              if (name.isNotEmpty) {
                final colorHex =
                    '#${selectedColor.value.value.toRadixString(16).substring(2)}';
                Provider.of<SubscriptionProvider>(
                  context,
                  listen: false,
                ).addLabel(name, colorHex).then((_) {
                  Navigator.of(context).pop();
                });
              }
            },
            child: const Text('Add'),
          ),
        ],
      ),
    );
  }

  // Build a color selection option
  Widget _buildColorOption(ValueNotifier<Color> selectedColor, Color color) {
    return ValueListenableBuilder<Color>(
      valueListenable: selectedColor,
      builder: (context, value, child) {
        final isSelected = value == color;
        return GestureDetector(
          onTap: () => selectedColor.value = color,
          child: Container(
            width: 30,
            height: 30,
            margin: const EdgeInsets.all(2),
            decoration: BoxDecoration(
              color: color,
              shape: BoxShape.circle,
              border: Border.all(
                color: isSelected ? Colors.white : Colors.transparent,
                width: 2,
              ),
              boxShadow: isSelected
                  ? [
                      BoxShadow(
                        color: Colors.black.withOpacity(0.3),
                        blurRadius: 4,
                      ),
                    ]
                  : null,
            ),
          ),
        );
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(
          _isEditMode ? 'Edit Subscription' : 'Add Subscription',
          style: const TextStyle(fontSize: 20, fontWeight: FontWeight.w600),
        ),
        elevation: 0,
      ),
      body: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Form(
            key: _formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                // Basic Information Section Title
                Padding(
                  padding: const EdgeInsets.only(left: 4, bottom: 12),
                  child: Text(
                    'Basic Information',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.w600,
                      color: Theme.of(context).colorScheme.primary,
                    ),
                  ),
                ),

                // Basic Information Card
                Card(
                  elevation: 2,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(16),
                  ),
                  child: Padding(
                    padding: const EdgeInsets.all(16.0),
                    child: Column(
                      children: [
                        TextFormField(
                          controller: _nameController,
                          decoration: const InputDecoration(
                            labelText: 'Subscription Name',
                            hintText: 'e.g., Netflix, Gym Membership',
                            prefixIcon: Icon(Icons.payment),
                          ),
                          textCapitalization: TextCapitalization.words,
                          validator: (value) {
                            if (value == null || value.trim().isEmpty) {
                              return 'Please enter a name';
                            }
                            return null;
                          },
                        ),
                        if (!_isEditMode ||
                            (_isEditMode && _isActiveSubscription))
                          const SizedBox(height: 16),
                        if (!_isEditMode ||
                            (_isEditMode && _isActiveSubscription))
                          Row(
                            children: [
                              Expanded(
                                flex: 7,
                                child: TextFormField(
                                  controller: _priceController,
                                  decoration: InputDecoration(
                                    labelText: 'Price',
                                    hintText: 'Enter the amount',
                                    prefixIcon: Icon(
                                      _selectedCurrency == Currency.USD.code
                                          ? Icons.attach_money
                                          : Icons.currency_exchange,
                                    ),
                                  ),
                                  enabled:
                                      !_isEditMode ||
                                      (_isEditMode && _isActiveSubscription),
                                  keyboardType:
                                      const TextInputType.numberWithOptions(
                                        decimal: true,
                                      ),
                                  inputFormatters: [
                                    FilteringTextInputFormatter.allow(
                                      RegExp(r'^\d+\.?\d{0,2}'),
                                    ),
                                  ],
                                  validator: (value) {
                                    if (value == null || value.isEmpty) {
                                      return 'Please enter a price';
                                    }
                                    try {
                                      final price = double.parse(value);
                                      if (price <= 0) {
                                        return 'Price must be greater than zero';
                                      }
                                    } catch (e) {
                                      return 'Please enter a valid number';
                                    }
                                    return null;
                                  },
                                ),
                              ),
                              const SizedBox(width: 8),
                              Expanded(
                                flex: 3,
                                child: DropdownButtonFormField<String>(
                                  value: _selectedCurrency,
                                  decoration: const InputDecoration(
                                    labelText: 'Currency',
                                    contentPadding: EdgeInsets.symmetric(
                                      horizontal: 10,
                                      vertical: 15,
                                    ),
                                  ),
                                  items: _currencies
                                      .map<DropdownMenuItem<String>>((
                                        String value,
                                      ) {
                                        return DropdownMenuItem<String>(
                                          value: value,
                                          child: Text(value),
                                        );
                                      })
                                      .toList(),
                                  onChanged:
                                      !_isEditMode ||
                                          (_isEditMode && _isActiveSubscription)
                                      ? (String? newValue) {
                                          if (newValue != null) {
                                            setState(() {
                                              _selectedCurrency = newValue;
                                            });
                                          }
                                        }
                                      : null,
                                ),
                              ),
                            ],
                          ),
                      ],
                    ),
                  ),
                ),

                if (!_isEditMode || (_isEditMode && _isActiveSubscription)) ...[
                  const SizedBox(height: 24),

                  // Subscription Details Section Title
                  Padding(
                    padding: const EdgeInsets.only(left: 4, bottom: 12),
                    child: Text(
                      'Subscription Details',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.w600,
                        color: Theme.of(context).colorScheme.primary,
                      ),
                    ),
                  ),

                  // Subscription Details Card
                  Card(
                    elevation: 2,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(16),
                    ),
                    child: Padding(
                      padding: const EdgeInsets.all(16.0),
                      child: Column(
                        children: [
                          DropdownButtonFormField<String>(
                            value: _selectedDuration,
                            decoration: const InputDecoration(
                              labelText: 'Subscription Recurrence',
                              prefixIcon: Icon(Icons.repeat),
                            ),
                            items: const [
                              DropdownMenuItem(
                                value: '1 month',
                                child: Text('Monthly'),
                              ),
                              DropdownMenuItem(
                                value: '3 months',
                                child: Text('Quarterly (3 months)'),
                              ),
                              DropdownMenuItem(
                                value: '6 months',
                                child: Text('Semi-Annual (6 months)'),
                              ),
                              DropdownMenuItem(
                                value: '12 months',
                                child: Text('Annual (12 months)'),
                              ),
                              DropdownMenuItem(
                                value: 'Custom',
                                child: Text('Custom'),
                              ),
                            ],
                            onChanged: _handleDurationChange,
                          ),
                          if (_selectedDuration == 'Custom')
                            Padding(
                              padding: const EdgeInsets.only(top: 16.0),
                              child: TextFormField(
                                controller: _customMonthsController,
                                decoration: const InputDecoration(
                                  labelText: 'Custom Months',
                                  hintText: 'Enter number of months',
                                  prefixIcon: Icon(Icons.calendar_today),
                                ),
                                keyboardType: TextInputType.number,
                                inputFormatters: [
                                  FilteringTextInputFormatter.digitsOnly,
                                ],
                                validator: (value) {
                                  if (_selectedDuration == 'Custom') {
                                    if (value == null || value.isEmpty) {
                                      return 'Please enter number of months';
                                    }
                                    final months = int.tryParse(value);
                                    if (months == null || months <= 0) {
                                      return 'Please enter a valid number';
                                    }
                                  }
                                  return null;
                                },
                                onChanged: _handleCustomMonthsChange,
                              ),
                            ),

                          const SizedBox(height: 16),

                          // Date picker
                          Container(
                            decoration: BoxDecoration(
                              borderRadius: BorderRadius.circular(12),
                              border: Border.all(
                                color: Theme.of(
                                  context,
                                ).dividerColor.withOpacity(0.5),
                              ),
                            ),
                            child: ListTile(
                              leading: Icon(
                                Icons.calendar_today,
                                color: Colors.blue.withOpacity(
                                  Theme.of(context).brightness ==
                                          Brightness.dark
                                      ? 0.8
                                      : 1.0,
                                ),
                              ),
                              title: const Text('Start Date'),
                              subtitle: Text(
                                '${_startDate.month}/${_startDate.day}/${_startDate.year}',
                              ),
                              onTap: () => _selectStartDate(context),
                            ),
                          ),
                          const SizedBox(height: 16),
                          // Toggle for recurrence count (only in add mode)
                          Row(
                            children: [
                              Switch(
                                value: _useRecurrenceCount,
                                onChanged: (value) {
                                  setState(() {
                                    _useRecurrenceCount = value;
                                    if (!value) {
                                      _recurrenceCountController.text = '1';
                                    }
                                  });
                                },
                              ),
                              const SizedBox(width: 8),
                              Text(
                                'Specify number of recurrences',
                                style: TextStyle(
                                  fontSize: 16,
                                  color: Theme.of(
                                    context,
                                  ).textTheme.bodyLarge?.color,
                                ),
                              ),
                            ],
                          ),

                          // Recurrence count input (shown only if toggle is on)
                          if (_useRecurrenceCount)
                            Padding(
                              padding: const EdgeInsets.only(
                                bottom: 8.0,
                                top: 8.0,
                              ),
                              child: TextFormField(
                                controller: _recurrenceCountController,
                                decoration: InputDecoration(
                                  labelText: 'Number of ${_getDurationLabel()}',
                                  hintText: 'Enter number of recurrences',
                                  prefixIcon: const Icon(Icons.repeat),
                                ),
                                keyboardType: TextInputType.number,
                                inputFormatters: [
                                  FilteringTextInputFormatter.digitsOnly,
                                ],
                                validator: (value) {
                                  if (_useRecurrenceCount) {
                                    if (value == null || value.isEmpty) {
                                      return 'Please enter number of recurrences';
                                    }
                                    final count = int.tryParse(value);
                                    if (count == null || count <= 0) {
                                      return 'Please enter a valid number';
                                    }
                                  }
                                  return null;
                                },
                              ),
                            ),
                        ],
                      ),
                    ),
                  ),
                ],

                const SizedBox(height: 24),

                // Family Section
                Consumer<FamilyMemberProvider>(
                  builder: (context, familyMemberProvider, _) {
                    final hasFamilyMembers = familyMemberProvider.hasFamilyMembers;

                    // If no family members, don't show the section
                    if (!hasFamilyMembers) {
                      return const SizedBox.shrink();
                    }

                    final familyMembers = familyMemberProvider.familyMembers;

                    return Card(
                      elevation: 2,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(16),
                      ),
                      child: Column(
                        children: [
                          // Collapsible header
                          InkWell(
                            onTap: () {
                              setState(() {
                                _isFamilyExpanded = !_isFamilyExpanded;
                              });
                            },
                            child: Padding(
                              padding: const EdgeInsets.all(16.0),
                              child: Row(
                                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                children: [
                                  Text(
                                    'Family',
                                    style: TextStyle(
                                      fontSize: 18,
                                      fontWeight: FontWeight.w600,
                                      color: Theme.of(context).colorScheme.primary,
                                    ),
                                  ),
                                  Icon(
                                    _isFamilyExpanded
                                        ? Icons.keyboard_arrow_up
                                        : Icons.keyboard_arrow_down,
                                    color: Theme.of(context).colorScheme.primary,
                                  ),
                                ],
                              ),
                            ),
                          ),

                          // Collapsible content
                          AnimatedContainer(
                            duration: const Duration(milliseconds: 200),
                            height: _isFamilyExpanded ? null : 0,
                            child: _isFamilyExpanded
                                ? Padding(
                                    padding: const EdgeInsets.fromLTRB(
                                      16.0,
                                      0,
                                      16.0,
                                      16.0,
                                    ),
                                    child: Column(
                                      crossAxisAlignment: CrossAxisAlignment.start,
                                      children: [
                                        const Divider(),
                                        const Text(
                                          'Link family members to this subscription',
                                          style: TextStyle(
                                            fontSize: 14,
                                            fontWeight: FontWeight.w400,
                                          ),
                                        ),
                                        const SizedBox(height: 16),

                                        // Who uses this subscription
                                        Column(
                                          crossAxisAlignment: CrossAxisAlignment.start,
                                          children: [
                                            const Text(
                                              'Who uses this subscription',
                                              style: TextStyle(
                                                fontSize: 16,
                                                fontWeight: FontWeight.w500,
                                              ),
                                            ),
                                            const SizedBox(height: 8),
                                            Wrap(
                                              spacing: 8.0,
                                              runSpacing: 8.0,
                                              children: _buildFamilyMemberChips(),
                                            ),
                                          ],
                                        ),

                                        const SizedBox(height: 16),

                                        // Who pays for this subscription
                                        DropdownButtonFormField<String?>(
                                          value: _payerFamilyMember?.id,
                                          decoration: const InputDecoration(
                                            labelText: 'Who pays for this subscription',
                                            prefixIcon: Icon(Icons.account_balance_wallet),
                                          ),
                                          items: [
                                            const DropdownMenuItem<String?>(
                                              value: null,
                                              child: Text('None'),
                                            ),
                                            const DropdownMenuItem<String?>(
                                              value: 'family',
                                              child: Text('Family (common account)'),
                                            ),
                                            ...familyMembers.where((member) => !member.isKid).map((member) {
                                              return DropdownMenuItem<String?>(
                                                value: member.id,
                                                child: Text(member.name),
                                              );
                                            }).toList(),
                                          ],
                                          onChanged: (String? value) {
                                            setState(() {
                                              if (value == 'family') {
                                                // Create a special FamilyMember object for "Family"
                                                _payerFamilyMember = FamilyMember(
                                                  id: 'family',
                                                  name: 'Family (common account)',
                                                );
                                              } else {
                                                _payerFamilyMember = value == null
                                                    ? null
                                                    : familyMemberProvider.getFamilyMemberById(value);
                                              }
                                            });
                                          },
                                        ),
                                      ],
                                    ),
                                  )
                                : const SizedBox.shrink(),
                          ),
                        ],
                      ),
                    );
                  },
                ),

                const SizedBox(height: 24),

                // Labels Section
                Card(
                  elevation: 2,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(16),
                  ),
                  child: Column(
                    children: [
                      // Collapsible header
                      InkWell(
                        onTap: () {
                          setState(() {
                            _isLabelsExpanded = !_isLabelsExpanded;
                          });
                        },
                        child: Padding(
                          padding: const EdgeInsets.all(16.0),
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              Text(
                                'Labels',
                                style: TextStyle(
                                  fontSize: 18,
                                  fontWeight: FontWeight.w600,
                                  color: Theme.of(context).colorScheme.primary,
                                ),
                              ),
                              Icon(
                                _isLabelsExpanded
                                    ? Icons.keyboard_arrow_up
                                    : Icons.keyboard_arrow_down,
                                color: Theme.of(context).colorScheme.primary,
                              ),
                            ],
                          ),
                        ),
                      ),

                      // Collapsible content
                      AnimatedContainer(
                        duration: const Duration(milliseconds: 200),
                        height: _isLabelsExpanded ? null : 0,
                        child: _isLabelsExpanded
                            ? Padding(
                                padding: const EdgeInsets.fromLTRB(
                                  16.0,
                                  0,
                                  16.0,
                                  16.0,
                                ),
                                child: Column(
                                  crossAxisAlignment: CrossAxisAlignment.start,
                                  children: [
                                    const Divider(),
                                    const Text(
                                      'Categorize your subscription',
                                      style: TextStyle(
                                        fontSize: 14,
                                        fontWeight: FontWeight.w400,
                                      ),
                                    ),
                                    const SizedBox(height: 16),
                                    Wrap(
                                      spacing: 8.0,
                                      runSpacing: 8.0,
                                      children: _buildLabelChips(),
                                    ),
                                    Align(
                                      alignment: Alignment.centerRight,
                                      child: TextButton.icon(
                                        icon: const Icon(Icons.add),
                                        label: const Text('Add Custom Label'),
                                        onPressed: _showAddLabelDialog,
                                      ),
                                    ),
                                  ],
                                ),
                              )
                            : const SizedBox.shrink(),
                      ),
                    ],
                  ),
                ),

                const SizedBox(height: 24),

                // Submit Button
                ElevatedButton(
                  onPressed: _submitForm,
                  style: ElevatedButton.styleFrom(
                    padding: const EdgeInsets.symmetric(vertical: 14),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                  child: Text(
                    _isEditMode ? 'Update Subscription' : 'Add Subscription',
                    style: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
