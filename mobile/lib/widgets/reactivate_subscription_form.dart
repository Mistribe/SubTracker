import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';

import '../models/subscription.dart';
import '../providers/subscription_provider.dart';

class ReactivateSubscriptionForm extends StatefulWidget {
  final Subscription subscription;

  const ReactivateSubscriptionForm({super.key, required this.subscription});

  @override
  State<ReactivateSubscriptionForm> createState() =>
      _ReactivateSubscriptionFormState();
}

class _ReactivateSubscriptionFormState
    extends State<ReactivateSubscriptionForm> {
  final _formKey = GlobalKey<FormState>();
  late final TextEditingController _priceController;
  late final TextEditingController _customMonthsController;
  late int _months;
  String _selectedDuration = '1 month'; // Default value
  DateTime _startDate = DateTime.now();
  DateTime? _endDate;
  late String _selectedCurrency;

  // List of common currencies
  final List<String> _currencies = ['USD', 'EUR', 'GBP', 'JPY', 'CAD', 'AUD', 'CHF', 'CNY', 'INR'];

  @override
  void initState() {
    super.initState();
    // Get the current subscription detail
    final currentDetail = widget.subscription.getLastPaymentDetail();

    // Initialize price with the last subscription price
    _priceController = TextEditingController(
      text: currentDetail.price.toString(),
    );
    _customMonthsController = TextEditingController();

    // Initialize months from subscription detail
    _months = currentDetail.months;
    _selectedCurrency = currentDetail.currency;

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
    super.dispose();
  }

  // Date picker for start date
  Future<void> _selectStartDate(BuildContext context) async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: _startDate,
      firstDate: widget.subscription.getLastPaymentDetail().endDate ?? DateTime.now(),
      lastDate: DateTime(2101),
    );
    if (picked != null && picked != _startDate) {
      setState(() {
        _startDate = picked;
      });
    }
  }

  // Date picker for end date
  Future<void> _selectEndDate(BuildContext context) async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: _endDate ?? _startDate.add(const Duration(days: 30)),
      firstDate: _startDate,
      lastDate: DateTime(2101),
    );
    if (picked != null && picked != _endDate) {
      setState(() {
        _endDate = picked;
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
    if (_formKey.currentState!.validate()) {
      final price = double.tryParse(_priceController.text) ?? 0.0;

      // For custom duration, get months from text field
      if (_selectedDuration == 'Custom') {
        final customMonths = int.tryParse(_customMonthsController.text);
        if (customMonths != null && customMonths > 0) {
          _months = customMonths;
        }
      }

      // Validate that if end date is provided, it's not shorter than the subscription recurrence
      if (_endDate != null) {
        // Calculate the minimum end date based on start date and months
        final minEndDate = DateTime(
          _startDate.year + (_months ~/ 12),
          _startDate.month + (_months % 12),
          _startDate.day,
        );

        if (_endDate!.isBefore(minEndDate)) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('End date must be at least $_months months after start date'),
              backgroundColor: Colors.red,
              duration: const Duration(seconds: 3),
            ),
          );
          return;
        }
      }

      try {
        // Show loading indicator
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Scheduling reactivation...'),
            duration: Duration(seconds: 1),
          ),
        );

        // Reactivate the subscription using the provider with all parameters
        await Provider.of<SubscriptionProvider>(
          context,
          listen: false,
        ).reactivatePayment(
          widget.subscription.id, 
          _startDate,
          price: price,
          months: _months,
          endDate: _endDate,
          currency: _selectedCurrency,
        );

        // Show success message
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Payment reactivation scheduled successfully'),
            duration: Duration(seconds: 2),
          ),
        );

        // Close the form
        Navigator.of(context).pop();
      } catch (e) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Error: ${e.toString()}'),
            backgroundColor: Colors.red,
            duration: const Duration(seconds: 3),
          ),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    // Get the keyboard height to adjust the bottom padding
    final mediaQuery = MediaQuery.of(context);
    final keyboardHeight = mediaQuery.viewInsets.bottom;

    return Padding(
      padding: EdgeInsets.only(
        top: 16,
        left: 16,
        right: 16,
        bottom: keyboardHeight + 16,
      ),
      child: SingleChildScrollView(
        child: Form(
          key: _formKey,
          child: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              const Text(
                'Reactivate Payment',
                style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 16),
              const Text(
                'Customize your subscription reactivation details below.',
              ),
              const SizedBox(height: 16),
              Row(
                children: [
                  Expanded(
                    flex: 7,
                    child: TextFormField(
                      controller: _priceController,
                      decoration: InputDecoration(
                        labelText: 'Price',
                        hintText: 'Enter the price',
                        prefixIcon: Icon(_selectedCurrency == 'USD' ? Icons.attach_money : Icons.currency_exchange),
                      ),
                      keyboardType: const TextInputType.numberWithOptions(
                        decimal: true,
                      ),
                      inputFormatters: [
                        FilteringTextInputFormatter.allow(RegExp(r'^\d+\.?\d{0,2}')),
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
                        contentPadding: EdgeInsets.symmetric(horizontal: 10, vertical: 15),
                      ),
                      items: _currencies.map<DropdownMenuItem<String>>((String value) {
                        return DropdownMenuItem<String>(
                          value: value,
                          child: Text(value),
                        );
                      }).toList(),
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
                  DropdownMenuItem(value: '1 month', child: Text('Monthly')),
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
                  DropdownMenuItem(value: 'Custom', child: Text('Custom')),
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
                    inputFormatters: [FilteringTextInputFormatter.digitsOnly],
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
              ListTile(
                leading: const Icon(Icons.calendar_today),
                title: const Text('Start Date'),
                subtitle: Text(
                  '${_startDate.month}/${_startDate.day}/${_startDate.year}',
                ),
                onTap: () => _selectStartDate(context),
              ),
              const SizedBox(height: 8),
              ListTile(
                leading: const Icon(Icons.event_busy),
                title: const Text('End Date (Optional)'),
                subtitle: Text(
                  _endDate != null
                      ? '${_endDate!.month}/${_endDate!.day}/${_endDate!.year}'
                      : 'No end date',
                ),
                trailing: _endDate != null
                    ? IconButton(
                        icon: const Icon(Icons.clear),
                        onPressed: () {
                          setState(() {
                            _endDate = null;
                          });
                        },
                      )
                    : null,
                onTap: () => _selectEndDate(context),
              ),
              const SizedBox(height: 24),
              ElevatedButton(
                onPressed: _submitForm,
                style: ElevatedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(vertical: 12),
                ),
                child: const Text('Reactivate', style: TextStyle(fontSize: 16)),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
