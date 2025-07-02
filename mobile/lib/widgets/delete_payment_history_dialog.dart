import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../models/subscription.dart';
import '../models/subscription_payment.dart';
import '../providers/subscription_provider.dart';

class DeletePaymentHistoryDialog {
  /// Shows a confirmation dialog for deleting a payment history entry
  /// Returns true if the payment history was successfully deleted, false otherwise
  static Future<bool> show({
    required BuildContext context,
    required Subscription subscription,
    required SubscriptionPayment paymentHistory,
  }) async {
    // Show confirmation dialog
    final shouldRemove = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Remove Payment History'),
        content: const Text(
          'Are you sure you want to remove this payment history? This action cannot be undone.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () => Navigator.of(context).pop(true),
            child: const Text('Remove'),
          ),
        ],
      ),
    ) ?? false;

    if (shouldRemove) {
      try {
        // Show loading indicator
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Removing payment history...'),
            duration: Duration(seconds: 1),
          ),
        );

        // Remove the payment history
        await Provider.of<SubscriptionProvider>(
          context,
          listen: false,
        ).removeSubscriptionPayment(
          subscription.id,
          paymentHistory.id,
        );

        // Show success message
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Payment history removed successfully'),
            duration: Duration(seconds: 2),
          ),
        );

        return true;
      } catch (e) {
        // Show error message
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              'Error removing payment history: ${e.toString()}',
            ),
            backgroundColor: Colors.red,
            duration: const Duration(seconds: 3),
          ),
        );
        return false;
      }
    }
    
    return false;
  }
}