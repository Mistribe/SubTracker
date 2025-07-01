import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../providers/payment_provider.dart';

class AddPaymentForm extends StatefulWidget {
  const AddPaymentForm({super.key});

  @override
  State<AddPaymentForm> createState() => _AddPaymentFormState();
}

class _AddPaymentFormState extends State<AddPaymentForm> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _priceController = TextEditingController();
  final _customMonthsController = TextEditingController();
  String _selectedDuration = '1 month'; // Default value
  DateTime _selectedDate = DateTime.now();

  @override
  void dispose() {
    _nameController.dispose();
    _priceController.dispose();
    _customMonthsController.dispose();
    super.dispose();
  }

  // Show date picker
  Future<void> _selectDate(BuildContext context) async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: _selectedDate,
      firstDate: DateTime(2000),
      lastDate: DateTime(2101),
    );
    if (picked != null && picked != _selectedDate) {
      setState(() {
        _selectedDate = picked;
      });
    }
  }

  Future<void> _submitForm() async {
    if (_formKey.currentState!.validate()) {
      final name = _nameController.text.trim();
      final price = double.parse(_priceController.text);

      // Determine the months value based on selected duration
      int months;
      if (_selectedDuration == '1 month') {
        months = 1;
      } else if (_selectedDuration == '3 months') {
        months = 3;
      } else if (_selectedDuration == '6 months') {
        months = 6;
      } else if (_selectedDuration == '12 months') {
        months = 12;
      } else {
        // Custom duration
        months = int.parse(_customMonthsController.text);
      }

      try {
        // Show loading indicator
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Adding payment...'),
            duration: Duration(seconds: 1),
          ),
        );

        // Add the payment using the provider
        await Provider.of<PaymentProvider>(
          context,
          listen: false,
        ).addPayment(name, price, months, _selectedDate);

        // Show success message
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Payment added successfully'),
            duration: Duration(seconds: 2),
          ),
        );

        // Close the form
        Navigator.of(context).pop();
      } catch (e) {
        // Show error message
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Error adding payment: ${e.toString()}'),
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
      child: Form(
        key: _formKey,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            const Text(
              'Add New Payment',
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
                hintText: 'Enter the amount',
                prefixIcon: Icon(Icons.attach_money),
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
            const SizedBox(height: 16),
            // Date picker
            ListTile(
              leading: const Icon(Icons.calendar_today),
              title: const Text('First Payment Date'),
              subtitle: Text(
                '${_selectedDate.month}/${_selectedDate.day}/${_selectedDate.year}',
              ),
              onTap: () => _selectDate(context),
            ),
            const SizedBox(height: 16),
            // Duration selector
            DropdownButtonFormField<String>(
              decoration: const InputDecoration(
                labelText: 'Payment Duration',
                hintText: 'Select payment duration',
                prefixIcon: Icon(Icons.calendar_today),
              ),
              value: _selectedDuration,
              items: [
                '1 month',
                '3 months',
                '6 months',
                '12 months',
                'Custom',
              ].map<DropdownMenuItem<String>>((String value) {
                return DropdownMenuItem<String>(
                  value: value,
                  child: Text(value),
                );
              }).toList(),
              onChanged: (String? newValue) {
                setState(() {
                  _selectedDuration = newValue!;
                });
              },
            ),
            // Show custom months input if 'Custom' is selected
            if (_selectedDuration == 'Custom')
              Padding(
                padding: const EdgeInsets.only(top: 16.0),
                child: TextFormField(
                  controller: _customMonthsController,
                  decoration: const InputDecoration(
                    labelText: 'Number of Months',
                    hintText: 'Enter the number of months',
                    prefixIcon: Icon(Icons.calendar_month),
                  ),
                  keyboardType: TextInputType.number,
                  inputFormatters: [
                    FilteringTextInputFormatter.digitsOnly,
                  ],
                  validator: (value) {
                    if (_selectedDuration == 'Custom') {
                      if (value == null || value.isEmpty) {
                        return 'Please enter the number of months';
                      }
                      try {
                        final months = int.parse(value);
                        if (months <= 0) {
                          return 'Number of months must be greater than zero';
                        }
                      } catch (e) {
                        return 'Please enter a valid number';
                      }
                    }
                    return null;
                  },
                ),
              ),
            const SizedBox(height: 24),
            ElevatedButton(
              onPressed: _submitForm,
              style: ElevatedButton.styleFrom(
                padding: const EdgeInsets.symmetric(vertical: 12),
              ),
              child: const Text('Add Payment', style: TextStyle(fontSize: 16)),
            ),
          ],
        ),
      ),
    );
  }
}
