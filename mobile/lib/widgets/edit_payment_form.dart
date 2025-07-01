import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../providers/payment_provider.dart';
import '../models/payment.dart';

class EditPaymentForm extends StatefulWidget {
  final Payment payment;

  const EditPaymentForm({super.key, required this.payment});

  @override
  State<EditPaymentForm> createState() => _EditPaymentFormState();
}

class _EditPaymentFormState extends State<EditPaymentForm> {
  final _formKey = GlobalKey<FormState>();
  late final TextEditingController _nameController;
  late final TextEditingController _priceController;
  late final TextEditingController _customMonthsController;
  late int _months;
  String _selectedDuration = '1 month'; // Default value
  late DateTime _startDate;
  DateTime? _endDate;

  @override
  void initState() {
    super.initState();
    // Initialize controllers with existing payment data
    _nameController = TextEditingController(text: widget.payment.name);

    // Get the current payment detail
    final currentDetail = widget.payment.getLastPaymentDetail();
    _priceController = TextEditingController(
      text: currentDetail.price.toString(),
    );
    _customMonthsController = TextEditingController();

    // Initialize months from payment detail
    _months = currentDetail.months;
    _startDate = currentDetail.startDate;
    _endDate = currentDetail.endDate;

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
    _nameController.dispose();
    _priceController.dispose();
    _customMonthsController.dispose();
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

  // Date picker for end date
  Future<void> _selectEndDate(BuildContext context) async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: _endDate ?? DateTime.now(),
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
      final name = _nameController.text.trim();
      final price = double.tryParse(_priceController.text) ?? 0.0;

      // For custom duration, get months from text field
      if (_selectedDuration == 'Custom') {
        final customMonths = int.tryParse(_customMonthsController.text);
        if (customMonths != null && customMonths > 0) {
          _months = customMonths;
        }
      }

      try {
        // Show loading indicator
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Updating payment...'),
            duration: Duration(seconds: 1),
          ),
        );

        // Update the payment name
        await Provider.of<PaymentProvider>(
          context,
          listen: false,
        ).updatePayment(widget.payment.id, name);

        // Update the payment detail
        final currentDetail = widget.payment.getLastPaymentDetail();
        await Provider.of<PaymentProvider>(
          context,
          listen: false,
        ).updatePaymentDetailEntry(
          widget.payment.id,
          currentDetail.id,
          price,
          _months,
          _startDate,
          endDate: _endDate,
        );

        // Show success message
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Payment updated successfully'),
            duration: Duration(seconds: 2),
          ),
        );

        // Close the form
        Navigator.of(context).pop();
      } catch (e) {
        // Show error message
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Error updating payment: ${e.toString()}'),
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

    // Get the current payment detail for display
    final currentDetail = widget.payment.getLastPaymentDetail();

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
                'Edit Payment',
                style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 16),
              TextFormField(
                controller: _nameController,
                decoration: const InputDecoration(
                  labelText: 'Payment Name',
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
              const SizedBox(height: 16),
              TextFormField(
                controller: _priceController,
                decoration: const InputDecoration(
                  labelText: 'Price',
                  hintText: 'Enter the price',
                  prefixIcon: Icon(Icons.attach_money),
                ),
                keyboardType: const TextInputType.numberWithOptions(decimal: true),
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
              const SizedBox(height: 16),
              DropdownButtonFormField<String>(
                value: _selectedDuration,
                decoration: const InputDecoration(
                  labelText: 'Payment Recurrence',
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
                child: const Text(
                  'Update Payment',
                  style: TextStyle(fontSize: 16),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
