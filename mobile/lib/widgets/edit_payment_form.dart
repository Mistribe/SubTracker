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

  @override
  void initState() {
    super.initState();
    // Initialize controllers with existing payment data
    _nameController = TextEditingController(text: widget.payment.name);

    // Get the current payment detail
    final currentDetail = widget.payment.getLastPaymentDetail();
    _priceController = TextEditingController(text: currentDetail.price.toString());
    _customMonthsController = TextEditingController();

    // Initialize months from payment detail
    _months = currentDetail.months;

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

  // Date picker method removed as payment date is now fixed

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
            content: Text('Updating payment...'),
            duration: Duration(seconds: 1),
          ),
        );

        // Update the payment name
        await Provider.of<PaymentProvider>(context, listen: false)
            .updatePayment(widget.payment.id, name);

        // Add a new payment detail entry with the updated price and months
        await Provider.of<PaymentProvider>(context, listen: false)
            .addPaymentDetailEntry(
              widget.payment.id,
              price,
              DateTime.now(),
              months: months,
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
              'Edit Payment',
              style: TextStyle(
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
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
              keyboardType: const TextInputType.numberWithOptions(decimal: true),
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
            // Payment date (non-editable)
            ListTile(
              leading: const Icon(Icons.calendar_today),
              title: const Text('Payment Date (Fixed)'),
              subtitle: Text('${widget.payment.getLastPaymentDetail().startDate.month}/${widget.payment.getLastPaymentDetail().startDate.day}/${widget.payment.getLastPaymentDetail().startDate.year}'),
              // Date is fixed and non-editable
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
              child: const Text(
                'Update Payment',
                style: TextStyle(fontSize: 16),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
