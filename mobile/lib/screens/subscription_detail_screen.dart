import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../models/subscription.dart';
import '../providers/subscription_provider.dart';
import '../widgets/cancel_subscription_form.dart';
import 'subscription_form_screen.dart';
import 'edit_subscription_payment_screen.dart';
import '../widgets/price_change_form.dart';
import '../widgets/reactivate_subscription_form.dart';
import '../widgets/delete_subscription_dialog.dart';
import '../widgets/delete_payment_history_dialog.dart';
import '../services/currency_converter.dart';

class PaymentDetailScreen extends StatefulWidget {
  final Subscription subscription;

  const PaymentDetailScreen({super.key, required this.subscription});

  @override
  State<PaymentDetailScreen> createState() => _PaymentDetailScreenState();
}

class _PaymentDetailScreenState extends State<PaymentDetailScreen> {
  late Subscription subscription;

  @override
  void initState() {
    super.initState();
    subscription = widget.subscription;
  }

  // Helper method to build summary item
  Widget _buildSummaryItem(
    BuildContext context,
    String title,
    String value,
    IconData icon,
    Color iconColor,
  ) {
    final isDark = Theme.of(context).brightness == Brightness.dark;

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8, horizontal: 4),
      child: Row(
        children: [
          Icon(
            icon,
            size: 18,
            color: iconColor.withOpacity(isDark ? 0.8 : 1.0),
          ),
          const SizedBox(width: 8),
          Text(
            title,
            style: TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w500,
              color: Theme.of(context).colorScheme.onSurface.withOpacity(0.7),
            ),
          ),
          const Spacer(),
          Text(
            value,
            style: TextStyle(
              fontSize: 16, 
              fontWeight: FontWeight.bold,
              color: Theme.of(context).colorScheme.primary,
            ),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    // Use Consumer to listen for changes in the PaymentProvider
    return Consumer<SubscriptionProvider>(
      builder: (context, subscriptionProvider, child) {
        // Get the updated subscription from the provider
        final updatedPayment = subscriptionProvider.subscriptions.firstWhere(
          (p) => p.id == subscription.id,
          orElse: () => subscription,
        );

        // Update the local subscription if it has changed
        if (updatedPayment != subscription) {
          subscription = updatedPayment;
        }

        // Helper function to build info rows
        Widget buildInfoRow(String label, String value, IconData icon) {
          final isDark = Theme.of(context).brightness == Brightness.dark;
          final iconColor =
              {
                Icons.attach_money: Colors.green,
                Icons.calendar_today: Colors.blue,
                Icons.calculate: Colors.purple,
                Icons.event_available: Colors.orange,
                Icons.monetization_on: Colors.amber,
                Icons.check_circle: Colors.green,
                Icons.stop_circle: Colors.red,
              }[icon] ??
              Theme.of(context).colorScheme.primary;

          return Container(
            margin: const EdgeInsets.only(bottom: 12.0),
            padding: const EdgeInsets.symmetric(vertical: 8, horizontal: 12),
            decoration: BoxDecoration(
              color: Theme.of(context).cardColor,
              borderRadius: BorderRadius.circular(12),
              border: Border.all(
                color: Theme.of(context).dividerColor.withOpacity(0.5),
                width: 1,
              ),
            ),
            child: Row(
              children: [
                Icon(
                  icon,
                  size: 20,
                  color: iconColor.withOpacity(isDark ? 0.8 : 1.0),
                ),
                const SizedBox(width: 12),
                Text(
                  label,
                  style: const TextStyle(
                    fontWeight: FontWeight.w500,
                    fontSize: 15,
                  ),
                ),
                const SizedBox(width: 8),
                Expanded(
                  child: Text(
                    value,
                    style: TextStyle(
                      fontSize: 15,
                      fontWeight: FontWeight.w600,
                      color: Theme.of(context).colorScheme.primary,
                    ),
                    textAlign: TextAlign.end,
                  ),
                ),
              ],
            ),
          );
        }

        // Get the current subscription detail
        final currentDetail = subscription.getLastPaymentDetail();

        // Determine subscription type based on months
        String subscriptionType;
        if (currentDetail.months == 1) {
          subscriptionType = 'Monthly';
        } else if (currentDetail.months == 3) {
          subscriptionType = 'Quarterly';
        } else if (currentDetail.months == 6) {
          subscriptionType = 'Semi-Annual';
        } else if (currentDetail.months == 12) {
          subscriptionType = 'Annual';
        } else {
          subscriptionType = 'Custom (${currentDetail.months} months)';
        }

        return Scaffold(
          appBar: AppBar(
            title: Text(
              subscription.name,
              style: const TextStyle(fontSize: 20, fontWeight: FontWeight.w600),
            ),
            elevation: 0,
            actions: [
              PopupMenuButton<String>(
                icon: const Icon(Icons.more_vert),
                onSelected: (value) async {
                  if (value == 'edit') {
                    Navigator.push(
                      context,
                      MaterialPageRoute(
                        builder: (context) =>
                            SubscriptionFormScreen(subscription: subscription),
                      ),
                    );
                  } else if (value == 'addPaymentHistory') {
                    showModalBottomSheet(
                      context: context,
                      isScrollControlled: true,
                      builder: (context) =>
                          PriceChangeForm(subscription: subscription),
                    );
                  } else if (value == 'stopPayment') {
                    showModalBottomSheet(
                      context: context,
                      isScrollControlled: true,
                      builder: (context) =>
                          StopSubscriptionForm(subscription: subscription),
                    );
                  } else if (value == 'reactivatePayment') {
                    showModalBottomSheet(
                      context: context,
                      isScrollControlled: true,
                      builder: (context) => ReactivateSubscriptionForm(
                        subscription: subscription,
                      ),
                    );
                  } else if (value == 'delete') {
                    await DeleteSubscriptionDialog.show(
                      context: context,
                      subscription: subscription,
                      navigateBack: true,
                    );
                  }
                },
                itemBuilder: (BuildContext context) => <PopupMenuEntry<String>>[
                  const PopupMenuItem<String>(
                    value: 'edit',
                    child: Row(
                      children: [
                        Icon(Icons.edit),
                        SizedBox(width: 8),
                        Text('Edit'),
                      ],
                    ),
                  ),
                  if (subscription.isActive && subscription.isStarted)
                    const PopupMenuItem<String>(
                      value: 'addPaymentHistory',
                      child: Row(
                        children: [
                          Icon(Icons.price_change),
                          SizedBox(width: 8),
                          Text('Price Change'),
                        ],
                      ),
                    ),
                  // Show stop subscription option if subscription is active
                  if (subscription.isActive && subscription.isStarted)
                    const PopupMenuItem<String>(
                      value: 'stopPayment',
                      child: Row(
                        children: [
                          Icon(Icons.stop_circle),
                          SizedBox(width: 8),
                          Text('Cancel'),
                        ],
                      ),
                    ),
                  // Show reactivate subscription option if subscription is not active
                  if (!subscription.isActive)
                    const PopupMenuItem<String>(
                      value: 'reactivatePayment',
                      child: Row(
                        children: [
                          Icon(Icons.play_circle),
                          SizedBox(width: 8),
                          Text('Reactivate'),
                        ],
                      ),
                    ),
                  const PopupMenuItem<String>(
                    value: 'delete',
                    child: Row(
                      children: [
                        Icon(Icons.delete_forever),
                        SizedBox(width: 8),
                        Text('Delete'),
                      ],
                    ),
                  ),
                ],
              ),
            ],
          ),
          body: SingleChildScrollView(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Title for the summary section
                Padding(
                  padding: const EdgeInsets.only(left: 4, bottom: 12),
                  child: Text(
                    'Summary',
                    style: TextStyle(
                      fontSize: 18,
                      fontWeight: FontWeight.w600,
                      color: Theme.of(context).colorScheme.primary,
                    ),
                  ),
                ),

                // Single summary card with all information
                Card(
                  elevation: 2,
                  shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
                  child: Padding(
                    padding: const EdgeInsets.all(16),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        // Price
                        Container(
                          margin: const EdgeInsets.only(bottom: 12.0),
                          child: Row(
                            children: [
                              Icon(
                                Icons.attach_money,
                                size: 20,
                                color: Colors.green,
                              ),
                              const SizedBox(width: 12),
                              const Text(
                                'Price',
                                style: TextStyle(
                                  fontWeight: FontWeight.w500,
                                  fontSize: 15,
                                ),
                              ),
                              const Spacer(),
                              Text(
                                CurrencyConverter.formatAmountWithCurrency(
                                  currentDetail.price,
                                  currentDetail.currency,
                                ),
                                style: TextStyle(
                                  fontSize: 16,
                                  fontWeight: FontWeight.w600,
                                  color: Theme.of(context).colorScheme.primary,
                                ),
                              ),
                            ],
                          ),
                        ),

                        const Divider(height: 1),

                        // Type
                        Container(
                          margin: const EdgeInsets.symmetric(vertical: 12.0),
                          child: Row(
                            children: [
                              Icon(
                                Icons.calendar_today,
                                size: 20,
                                color: Colors.blue,
                              ),
                              const SizedBox(width: 12),
                              const Text(
                                'Type',
                                style: TextStyle(
                                  fontWeight: FontWeight.w500,
                                  fontSize: 15,
                                ),
                              ),
                              const Spacer(),
                              Text(
                                subscriptionType,
                                style: TextStyle(
                                  fontSize: 16,
                                  fontWeight: FontWeight.w600,
                                  color: Theme.of(context).colorScheme.primary,
                                ),
                              ),
                            ],
                          ),
                        ),

                        const Divider(height: 1),

                        // Status
                        Container(
                          margin: const EdgeInsets.symmetric(vertical: 12.0),
                          child: Row(
                            children: [
                              Icon(
                                subscription.isActive ? Icons.check_circle : Icons.stop_circle,
                                size: 20,
                                color: subscription.isActive ? Colors.green : Colors.red,
                              ),
                              const SizedBox(width: 12),
                              const Text(
                                'Status',
                                style: TextStyle(
                                  fontWeight: FontWeight.w500,
                                  fontSize: 15,
                                ),
                              ),
                              const Spacer(),
                              Text(
                                subscription.isActive ? 'Active' : 'Stopped',
                                style: TextStyle(
                                  fontSize: 16,
                                  fontWeight: FontWeight.w600,
                                  color: Theme.of(context).colorScheme.primary,
                                ),
                              ),
                            ],
                          ),
                        ),

                        // Monthly cost (if applicable)
                        if (currentDetail.months > 1) ...[
                          const Divider(height: 1),
                          Container(
                            margin: const EdgeInsets.symmetric(vertical: 12.0),
                            child: Row(
                              children: [
                                Icon(
                                  Icons.calculate,
                                  size: 20,
                                  color: Colors.purple,
                                ),
                                const SizedBox(width: 12),
                                const Text(
                                  'Monthly',
                                  style: TextStyle(
                                    fontWeight: FontWeight.w500,
                                    fontSize: 15,
                                  ),
                                ),
                                const Spacer(),
                                Text(
                                  CurrencyConverter.formatAmountWithCurrency(
                                    subscription.monthlyCost,
                                    currentDetail.currency,
                                  ),
                                  style: TextStyle(
                                    fontSize: 16,
                                    fontWeight: FontWeight.w600,
                                    color: Theme.of(context).colorScheme.primary,
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ],

                        // Next payment (if active)
                        if (subscription.isActive) ...[
                          const Divider(height: 1),
                          Container(
                            margin: const EdgeInsets.symmetric(vertical: 12.0),
                            child: Row(
                              children: [
                                Icon(
                                  Icons.event_available,
                                  size: 20,
                                  color: Colors.orange,
                                ),
                                const SizedBox(width: 12),
                                const Text(
                                  'Next Payment',
                                  style: TextStyle(
                                    fontWeight: FontWeight.w500,
                                    fontSize: 15,
                                  ),
                                ),
                                const Spacer(),
                                Text(
                                  subscription.formattedNextPaymentDate,
                                  style: TextStyle(
                                    fontSize: 16,
                                    fontWeight: FontWeight.w600,
                                    color: Theme.of(context).colorScheme.primary,
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ],

                        const Divider(height: 1),

                        // Total spent
                        Container(
                          margin: const EdgeInsets.symmetric(vertical: 12.0),
                          child: Row(
                            children: [
                              Icon(
                                Icons.monetization_on,
                                size: 20,
                                color: Colors.amber,
                              ),
                              const SizedBox(width: 12),
                              const Text(
                                'Total Spent',
                                style: TextStyle(
                                  fontWeight: FontWeight.w500,
                                  fontSize: 15,
                                ),
                              ),
                              const Spacer(),
                              Text(
                                subscription.formattedTotalAmountSpent,
                                style: TextStyle(
                                  fontSize: 16,
                                  fontWeight: FontWeight.w600,
                                  color: Theme.of(context).colorScheme.primary,
                                ),
                              ),
                            ],
                          ),
                        ),

                        // End date (if not active)
                        if (!subscription.isActive && currentDetail.endDate != null) ...[
                          const Divider(height: 1),
                          Container(
                            margin: const EdgeInsets.symmetric(vertical: 12.0),
                            child: Row(
                              children: [
                                Icon(
                                  Icons.stop_circle,
                                  size: 20,
                                  color: Colors.red,
                                ),
                                const SizedBox(width: 12),
                                const Text(
                                  'Stop Date',
                                  style: TextStyle(
                                    fontWeight: FontWeight.w500,
                                    fontSize: 15,
                                  ),
                                ),
                                const Spacer(),
                                Text(
                                  '${currentDetail.endDate!.month}/${currentDetail.endDate!.day}/${currentDetail.endDate!.year}',
                                  style: TextStyle(
                                    fontSize: 16,
                                    fontWeight: FontWeight.w600,
                                    color: Theme.of(context).colorScheme.primary,
                                  ),
                                ),
                              ],
                            ),
                          ),
                        ],
                      ],
                    ),
                  ),
                ),

                // Show labels if there are any
                if (subscription.labels.isNotEmpty)
                  Padding(
                    padding: const EdgeInsets.only(top: 12),
                    child: Card(
                      elevation: 2,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(16),
                      ),
                      child: Padding(
                        padding: const EdgeInsets.all(12),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Row(
                              children: [
                                Icon(
                                  Icons.label,
                                  size: 16,
                                  color: Theme.of(context).primaryColor,
                                ),
                                const SizedBox(width: 6),
                                Text(
                                  'Labels',
                                  style: TextStyle(
                                    fontSize: 14,
                                    fontWeight: FontWeight.w500,
                                    color: Theme.of(
                                      context,
                                    ).colorScheme.onSurface.withOpacity(0.7),
                                  ),
                                ),
                              ],
                            ),
                            const SizedBox(height: 8),
                            Wrap(
                              spacing: 4,
                              runSpacing: 4,
                              children: subscription.labels.map((label) {
                                return Chip(
                                  label: Text(
                                    label.name,
                                    style: const TextStyle(
                                      fontSize: 12,
                                      color: Colors.white,
                                      fontWeight: FontWeight.w500,
                                    ),
                                  ),
                                  backgroundColor: Color(
                                    int.parse(
                                          label.color.substring(1, 7),
                                          radix: 16,
                                        ) +
                                        0xFF000000,
                                  ),
                                  materialTapTargetSize:
                                      MaterialTapTargetSize.shrinkWrap,
                                  padding: const EdgeInsets.symmetric(
                                    horizontal: 8,
                                    vertical: 0,
                                  ),
                                  shape: RoundedRectangleBorder(
                                    borderRadius: BorderRadius.circular(16),
                                  ),
                                );
                              }).toList(),
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),

                const SizedBox(height: 24),

                // Payment History section
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 4),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Text(
                        'History',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.w600,
                          color: Theme.of(context).colorScheme.primary,
                        ),
                      ),
                      Text(
                        '${subscription.subscriptionPayments.length} entries',
                        style: TextStyle(
                          fontSize: 14,
                          color: Theme.of(
                            context,
                          ).colorScheme.onSurface.withOpacity(0.6),
                        ),
                      ),
                    ],
                  ),
                ),
                const SizedBox(height: 12),

                if (subscription.subscriptionPayments.isEmpty)
                  Card(
                    elevation: 2,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(16),
                    ),
                    child: const Padding(
                      padding: EdgeInsets.all(16.0),
                      child: Text(
                        'No subscription history recorded yet.',
                        style: TextStyle(fontSize: 16),
                      ),
                    ),
                  )
                else
                  Card(
                    elevation: 2,
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(16),
                    ),
                    child: ListView.separated(
                      shrinkWrap: true,
                      physics: const NeverScrollableScrollPhysics(),
                      itemCount: subscription.subscriptionPayments.length,
                      separatorBuilder: (context, index) => const Divider(),
                      itemBuilder: (context, index) {
                        final history =
                            subscription.subscriptionPayments[subscription
                                    .subscriptionPayments
                                    .length -
                                1 -
                                index];
                        return ListTile(
                          leading: Icon(
                            !history.isStarted
                                ? Icons.schedule
                                : history.isActive
                                ? Icons.check_circle
                                : Icons.stop_circle,
                            color: !history.isStarted
                                ? Colors.orange
                                : history.isActive
                                ? Colors.green
                                : Colors.red,
                          ),
                          title: Text(
                            CurrencyConverter.formatAmountWithCurrency(
                              history.price,
                              history.currency,
                            ),
                          ),
                          subtitle: Text(
                            !history.isStarted
                                ? 'Start at ${history.startDate.month}/${history.startDate.day}/${history.startDate.year}'
                                : 'From ${history.startDate.month}/${history.startDate.day}/${history.startDate.year} to ${history.endDate == null ? "now" : "${history.endDate!.month}/${history.endDate!.day}/${history.endDate!.year}"}',
                          ),
                          trailing: PopupMenuButton<String>(
                            icon: const Icon(Icons.more_vert),
                            onSelected: (value) async {
                              if (value == 'edit') {
                                Navigator.push(
                                  context,
                                  MaterialPageRoute(
                                    builder: (context) =>
                                        EditSubscriptionPaymentScreen(
                                          subscription: subscription,
                                          paymentHistory: history,
                                        ),
                                  ),
                                );
                              } else if (value == 'remove') {
                                await DeletePaymentHistoryDialog.show(
                                  context: context,
                                  subscription: subscription,
                                  paymentHistory: history,
                                );
                              }
                            },
                            itemBuilder: (BuildContext context) =>
                                <PopupMenuEntry<String>>[
                                  const PopupMenuItem<String>(
                                    value: 'edit',
                                    child: Row(
                                      children: [
                                        Icon(Icons.edit),
                                        SizedBox(width: 8),
                                        Text('Edit'),
                                      ],
                                    ),
                                  ),
                                  if (subscription.subscriptionPayments.length >
                                      1)
                                    const PopupMenuItem<String>(
                                      value: 'remove',
                                      child: Row(
                                        children: [
                                          Icon(Icons.delete),
                                          SizedBox(width: 8),
                                          Text('Remove'),
                                        ],
                                      ),
                                    ),
                                ],
                          ),
                        );
                      },
                    ),
                  ),

                const SizedBox(height: 24),
              ],
            ),
          ),
        );
      },
    );
  }
}
