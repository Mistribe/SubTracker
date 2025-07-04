import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../models/subscription.dart';
import '../providers/subscription_provider.dart';

class StopSubscriptionForm extends StatefulWidget {
  final Subscription subscription;

  const StopSubscriptionForm({super.key, required this.subscription});

  @override
  State<StopSubscriptionForm> createState() => _CancelSubscriptionFormState();
}

class _CancelSubscriptionFormState extends State<StopSubscriptionForm> {
  final _formKey = GlobalKey<FormState>();
  DateTime _selectedDate = DateTime.now();
  bool _useLastPaymentDate = true;

  Future<void> _submitForm() async {
    final messenger = ScaffoldMessenger.of(context);
    final navigator = Navigator.of(context);
    final subscriptionProvider = Provider.of<SubscriptionProvider>(
      context,
      listen: false,
    );

    if (_formKey.currentState!.validate()) {
      try {
        // Show loading indicator
        messenger.showSnackBar(
          const SnackBar(
            content: Text('Canceling subscription...'),
            duration: Duration(seconds: 1),
          ),
        );

        // Update the subscription detail
        await subscriptionProvider.cancelCurrentSubscription(
          widget.subscription.id,
          stopDate: _selectedDate,
        );

        // Show success message
        messenger.showSnackBar(
          const SnackBar(
            content: Text('Subscription successfully cancelled'),
            duration: Duration(seconds: 2),
          ),
        );

        // Close the form
        navigator.pop();
      } catch (e) {
        // Show error message
        messenger.showSnackBar(
          SnackBar(
            content: Text('Error canceling subscription: ${e.toString()}'),
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
              Text(
                'Cancelling ${widget.subscription.name}',
                style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 16),
              const Text(
                'Are you sure you want to stop this subscription? '
                'A stopped subscription will not be counted in active or total subscriptions.',
              ),
              const SizedBox(height: 16),
              SwitchListTile(
                title: const Text('Use last subscription date'),
                subtitle: const Text(
                  'If enabled, the subscription will be stopped at the last subscription date',
                ),
                value: _useLastPaymentDate,
                onChanged: (value) {
                  setState(() {
                    _useLastPaymentDate = value;
                  });
                },
              ),
              if (!_useLastPaymentDate)
                ListTile(
                  leading: const Icon(Icons.calendar_today),
                  title: const Text('Stop Date'),
                  subtitle: Text(
                    '${_selectedDate.month}/${_selectedDate.day}/${_selectedDate.year}',
                  ),
                  onTap: () async {
                    final DateTime? picked = await showDatePicker(
                      context: context,
                      initialDate: _selectedDate,
                      firstDate: widget.subscription
                          .getLastPaymentDetail()
                          .startDate,
                      lastDate: DateTime(DateTime.now().year + 100),
                    );
                    if (picked != null && picked != _selectedDate) {
                      setState(() {
                        _selectedDate = picked;
                      });
                    }
                  },
                ),
              const SizedBox(height: 24),
              ElevatedButton(
                onPressed: _submitForm,
                style: ElevatedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(vertical: 12),
                ),
                child: const Text(
                  'Cancel subscription',
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
