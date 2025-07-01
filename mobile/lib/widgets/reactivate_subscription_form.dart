import 'package:flutter/material.dart';
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
  DateTime _selectedDate = DateTime.now();

  Future<void> _selectReactivationDate(BuildContext context) async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: _selectedDate,
      firstDate:
          widget.subscription.getLastPaymentDetail().endDate ?? DateTime.now(),
      lastDate: DateTime(2101),
    );
    if (picked != null && picked != _selectedDate) {
      setState(() {
        _selectedDate = picked;
      });
    }
  }

  Future<void> _submitForm() async {
    try {
      // Show loading indicator
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Scheduling reactivation...'),
          duration: Duration(seconds: 1),
        ),
      );

      // Reactivate the subscription using the provider
      await Provider.of<SubscriptionProvider>(
        context,
        listen: false,
      ).reactivatePayment(widget.subscription.id, _selectedDate);

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
              'Select a date to reactivate this subscription. '
              'The subscription will be reactivated on this date.',
            ),
            const SizedBox(height: 16),
            ListTile(
              leading: const Icon(Icons.calendar_today),
              title: const Text('Reactivation Date'),
              subtitle: Text(
                '${_selectedDate.month}/${_selectedDate.day}/${_selectedDate.year}',
              ),
              onTap: () => _selectReactivationDate(context),
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
    );
  }
}
