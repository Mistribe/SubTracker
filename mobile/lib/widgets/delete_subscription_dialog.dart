import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../models/subscription.dart';
import '../providers/subscription_provider.dart';

class DeleteSubscriptionDialog {
  /// Shows a confirmation dialog for deleting a subscription
  /// Returns true if the subscription was successfully deleted, false otherwise
  static Future<bool> show({
    required BuildContext context,
    required Subscription subscription,
    bool navigateBack = false,
  }) async {
    // Show confirmation dialog
    final shouldDelete = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Delete Subscription'),
        content: const Text(
          'Are you sure you want to delete this subscription? This action is irreversible and cannot be undone.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('Cancel'),
          ),
          TextButton(
            onPressed: () => Navigator.of(context).pop(true),
            child: const Text('Delete'),
          ),
        ],
      ),
    ) ?? false;

    if (shouldDelete) {
      try {
        // Show loading indicator
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Removing subscription...'),
            duration: Duration(seconds: 1),
          ),
        );

        // Remove the subscription
        await Provider.of<SubscriptionProvider>(
          context,
          listen: false,
        ).removePayment(
          subscription.id,
        );

        // Show success message
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('${subscription.name} removed'),
            duration: const Duration(seconds: 2),
          ),
        );

        // Navigate back if requested (for detail screen)
        if (navigateBack) {
          Navigator.of(context).pop();
        }

        return true;
      } catch (e) {
        // Show error message
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text(
              'Error removing subscription: ${e.toString()}',
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