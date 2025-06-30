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
  late bool _isAnnual;
  @override
  void initState() {
    super.initState();
    // Initialize controllers with existing payment data
    _nameController = TextEditingController(text: widget.payment.name);
    _priceController = TextEditingController(text: widget.payment.price.toString());
    _isAnnual = widget.payment.isAnnual;
  }

  @override
  void dispose() {
    _nameController.dispose();
    _priceController.dispose();
    super.dispose();
  }

  // Date picker method removed as payment date is now fixed

  Future<void> _submitForm() async {
    if (_formKey.currentState!.validate()) {
      final name = _nameController.text.trim();
      final price = double.parse(_priceController.text);

      try {
        // Show loading indicator
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Updating payment...'),
            duration: Duration(seconds: 1),
          ),
        );

        // Create updated payment object
        // Note: We're using the original payment date to ensure it remains fixed
        final updatedPayment = widget.payment.copyWith(
          name: name,
          price: price,
          isAnnual: _isAnnual,
          // Original payment date is preserved (not editable)
        );

        // Update the payment using the provider
        await Provider.of<PaymentProvider>(context, listen: false)
            .updatePayment(updatedPayment);

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
              subtitle: Text('${widget.payment.paymentDate.month}/${widget.payment.paymentDate.day}/${widget.payment.paymentDate.year}'),
              // Date is fixed and non-editable
            ),
            const SizedBox(height: 16),
            SwitchListTile(
              title: const Text('Annual Payment'),
              subtitle: const Text('Toggle if this is paid yearly instead of monthly'),
              value: _isAnnual,
              onChanged: (value) {
                setState(() {
                  _isAnnual = value;
                });
              },
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
