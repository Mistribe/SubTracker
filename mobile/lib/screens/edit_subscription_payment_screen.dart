import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../models/subscription.dart';
import '../models/subscription_payment.dart';
import '../models/currency.dart';
import '../providers/subscription_provider.dart';

class EditSubscriptionPaymentScreen extends StatefulWidget {
  final Subscription subscription;
  final SubscriptionPayment paymentHistory;

  const EditSubscriptionPaymentScreen({
    super.key,
    required this.subscription,
    required this.paymentHistory,
  });

  @override
  State<EditSubscriptionPaymentScreen> createState() =>
      _EditSubscriptionPaymentScreenState();
}

class _EditSubscriptionPaymentScreenState
    extends State<EditSubscriptionPaymentScreen> {
  final _formKey = GlobalKey<FormState>();
  late final TextEditingController _priceController;
  late final TextEditingController _customMonthsController;
  late final TextEditingController _recurrenceCountController;
  late int _months;
  String _selectedDuration = '1 month'; // Default value
  late DateTime _startDate;
  DateTime? _endDate;
  bool _useRecurrenceCount = false;
  late String _selectedCurrency;

  // Get currencies from the Currency enum
  final List<String> _currencies = Currency.codes;

  // Calculate recurrence count based on start date, end date, and months
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

  @override
  void initState() {
    super.initState();
    // Initialize controllers with existing payment history data
    _priceController = TextEditingController(
      text: widget.paymentHistory.price.toString(),
    );
    _customMonthsController = TextEditingController();
    _recurrenceCountController = TextEditingController(text: '1');

    // Initialize months from payment history
    _months = widget.paymentHistory.months;
    _startDate = widget.paymentHistory.startDate;
    _endDate = widget.paymentHistory.endDate;
    _selectedCurrency = widget.paymentHistory.currency;

    // Always enable recurrence count
    if (_endDate != null) {
      // If end date is provided, calculate recurrence count
      int recurrenceCount = _calculateRecurrenceCount(
        _startDate,
        _endDate!,
        _months,
      );
      _recurrenceCountController.text = recurrenceCount.toString();
    } else {
      // Default to 1 recurrence
      _recurrenceCountController.text = '1';
    }

    // Always enable recurrence count as it's the only way to set end date
    _useRecurrenceCount = true;

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
      _customMonthsController.text = _months.toString();
    }
  }

  @override
  void dispose() {
    _priceController.dispose();
    _customMonthsController.dispose();
    _recurrenceCountController.dispose();
    super.dispose();
  }

  // Date picker for start date
  Future<void> _selectStartDate(BuildContext context) async {
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

  // Handle custom months change
  void _handleCustomMonthsChange(String value) {
    if (value.isEmpty) return;

    final parsedValue = int.tryParse(value);
    if (parsedValue != null && parsedValue > 0) {
      setState(() {
        _months = parsedValue;
      });
    }
  }

  Future<void> _submitForm() async {
    final messenger = ScaffoldMessenger.of(context);
    final navigator = Navigator.of(context);
    final subscriptionProvider = Provider.of<SubscriptionProvider>(
      context,
      listen: false,
    );

    if (_formKey.currentState!.validate()) {
      final price = double.tryParse(_priceController.text) ?? 0.0;

      // For custom duration, get months from text field
      if (_selectedDuration == 'Custom') {
        final customMonths = int.tryParse(_customMonthsController.text);
        if (customMonths != null && customMonths > 0) {
          _months = customMonths;
        }
      }

      // Calculate the end date based on recurrence count
      if (_useRecurrenceCount) {
        final recurrenceCount =
            int.tryParse(_recurrenceCountController.text) ?? 1;
        if (recurrenceCount > 0) {
          _endDate = _calculateEndDate();
        }
      } else {
        // If recurrence count is not used, set end date to null
        _endDate = null;
      }

      try {
        // Show loading indicator
        messenger.showSnackBar(
          const SnackBar(
            content: Text('Updating payment history...'),
            duration: Duration(seconds: 1),
          ),
        );

        // Update the payment history
        await subscriptionProvider.updateSubscriptionPayment(
          widget.subscription.id,
          widget.paymentHistory.id,
          price,
          _months,
          _startDate,
          endDate: _endDate,
          currency: _selectedCurrency,
        );

        // Show success message
        messenger.showSnackBar(
          const SnackBar(
            content: Text('Payment history updated successfully'),
            duration: Duration(seconds: 2),
          ),
        );

        // Close the form
        navigator.pop();
      } catch (e) {
        // Show error message
        messenger.showSnackBar(
          SnackBar(
            content: Text('Error updating payment history: ${e.toString()}'),
            backgroundColor: Colors.red,
            duration: const Duration(seconds: 3),
          ),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text(
          'Edit Payment History',
          style: TextStyle(fontSize: 20, fontWeight: FontWeight.w600),
        ),
        elevation: 0,
        leading: IconButton(
          icon: const Icon(Icons.close),
          onPressed: () => Navigator.of(context).pop(),
        ),
        actions: [
          IconButton(
            icon: const Icon(Icons.check),
            onPressed: _submitForm,
            tooltip: 'Update Payment History',
          ),
        ],
      ),
      body: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(16.0),
          child: Form(
            key: _formKey,
            child: Column(
              mainAxisSize: MainAxisSize.min,
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                // Payment Details Section Title
                Padding(
                  padding: const EdgeInsets.only(left: 4, bottom: 12),
                  child: Text(
                    'Payment Details',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.w600,
                      color: Theme.of(context).colorScheme.primary,
                    ),
                  ),
                ),

                // Payment Details Card
                Card(
                  elevation: 2,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(16),
                  ),
                  child: Padding(
                    padding: const EdgeInsets.all(16.0),
                    child: Column(
                      children: [
                        Row(
                          children: [
                            Expanded(
                              flex: 7,
                              child: TextFormField(
                                controller: _priceController,
                                decoration: InputDecoration(
                                  labelText: 'Price',
                                  hintText: 'Enter the price',
                                  prefixIcon: Icon(
                                    _selectedCurrency == Currency.USD.code
                                        ? Icons.attach_money
                                        : Icons.currency_exchange,
                                  ),
                                ),
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
                                  final price = double.tryParse(value);
                                  if (price == null || price <= 0) {
                                    return 'Please enter a valid price';
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
                                onChanged: (String? newValue) {
                                  if (newValue != null) {
                                    setState(() {
                                      _selectedCurrency = newValue;
                                    });
                                  }
                                },
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 16),
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

                        // Date pickers with modern styling
                        Container(
                          decoration: BoxDecoration(
                            borderRadius: BorderRadius.circular(12),
                            border: Border.all(
                              color: Theme.of(
                                context,
                              ).dividerColor.withValues(alpha: 0.5),
                            ),
                          ),
                          child: ListTile(
                            leading: Icon(
                              Icons.calendar_today,
                              color: Colors.blue.withValues(
                                alpha:
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

                        // Toggle for recurrence count
                        Row(
                          children: [
                            Switch(
                              value: _useRecurrenceCount,
                              onChanged: (value) {
                                setState(() {
                                  _useRecurrenceCount = value;
                                  if (!value) {
                                    _recurrenceCountController.text = '1';
                                    _endDate = null;
                                  } else {
                                    // Calculate end date when toggle is turned on
                                    final recurrenceCount =
                                        int.tryParse(
                                          _recurrenceCountController.text,
                                        ) ??
                                        1;
                                    if (recurrenceCount > 0) {
                                      _endDate = _calculateEndDate();
                                    }
                                  }
                                });
                              },
                            ),
                            const SizedBox(width: 8),
                            Expanded(
                              child: Text(
                                'Set end date by specifying number of recurrences',
                                style: TextStyle(
                                  fontSize: 16,
                                  color: Theme.of(
                                    context,
                                  ).textTheme.bodyLarge?.color,
                                ),
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
            ),
          ),
        ),
      ),
    );
  }
}
