import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:subscription_tracker/widgets/price_change_form.dart';
import '../providers/subscription_provider.dart';
import '../models/subscription.dart';
import '../models/currency.dart';
import '../services/currency_converter.dart';
import '../screens/subscription_detail_screen.dart';
import '../screens/edit_subscription_screen.dart';
import '../widgets/delete_subscription_dialog.dart';

class SubscriptionList extends StatelessWidget {
  const SubscriptionList({super.key});

  @override
  Widget build(BuildContext context) {
    final subscriptionProvider = Provider.of<SubscriptionProvider>(context);
    final subscriptions = subscriptionProvider.filteredSubscriptions;

    if (subscriptions.isEmpty) {
      // Check if we have subscriptions but they're filtered out
      final hasSubscriptions = subscriptionProvider.subscriptions.isNotEmpty;
      final isFiltered = !subscriptionProvider.showInactiveSubscriptions;

      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              hasSubscriptions && isFiltered 
                ? Icons.filter_alt_outlined 
                : Icons.subscriptions_outlined,
              size: 64,
              color: Theme.of(context).colorScheme.primary.withOpacity(0.5),
            ),
            const SizedBox(height: 16),
            Text(
              hasSubscriptions && isFiltered
                ? 'No active subscriptions'
                : 'No subscriptions yet',
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                color: Theme.of(context).colorScheme.onSurface,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              hasSubscriptions && isFiltered
                ? 'Toggle the filter to show inactive subscriptions'
                : 'Add your first subscription!',
              style: TextStyle(
                fontSize: 16,
                color: Theme.of(context).colorScheme.onSurface.withOpacity(0.7),
              ),
            ),
          ],
        ),
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
      itemCount: subscriptions.length,
      itemBuilder: (context, index) {
        final subscription = subscriptions[index];
        return SubscriptionCard(subscription: subscription);
      },
    );
  }
}

class SubscriptionCard extends StatelessWidget {
  final Subscription subscription;

  const SubscriptionCard({super.key, required this.subscription});

  Widget _buildFrequencyBadge(Subscription subscription) {
    if (!subscription.isActive || subscription.subscriptionPayments.isEmpty) {
      return const SizedBox.shrink();
    }

    final lastPayment = subscription.getLastPaymentDetail();
    final months = lastPayment.months;

    String frequencyText;
    Color badgeColor;

    if (months == 1) {
      frequencyText = 'Monthly';
      badgeColor = Colors.blue;
    } else if (months == 3) {
      frequencyText = 'Quarterly';
      badgeColor = Colors.green;
    } else if (months == 6) {
      frequencyText = 'Bi-annually';
      badgeColor = Colors.orange;
    } else if (months == 12) {
      frequencyText = 'Annually';
      badgeColor = Colors.purple;
    } else {
      frequencyText = '$months months';
      badgeColor = Colors.teal;
    }

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: badgeColor.withOpacity(0.2),
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: badgeColor, width: 1),
      ),
      child: Text(
        frequencyText,
        style: TextStyle(
          fontSize: 12,
          fontWeight: FontWeight.bold,
          color: badgeColor,
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final subscriptionProvider = Provider.of<SubscriptionProvider>(
      context,
      listen: false,
    );
    final isDark = Theme.of(context).brightness == Brightness.dark;

    // Get currencies
    final defaultCurrency = subscriptionProvider.defaultCurrencyEnum;
    final subscriptionCurrency = Currency.fromCode(
      subscription.subscriptionPayments.isNotEmpty
          ? subscription.getLastPaymentDetail().currency
          : Currency.defaultCode
    );

    // Convert costs to default currency
    final monthlyCostInDefaultCurrency = CurrencyConverter.convert(
      subscription.monthlyCost,
      subscriptionCurrency.code,
      defaultCurrency.code,
    );
    final annualCostInDefaultCurrency = monthlyCostInDefaultCurrency * 12;
    final totalSpentInDefaultCurrency = CurrencyConverter.convert(
      subscription.totalAmountSpent,
      subscriptionCurrency.code,
      defaultCurrency.code,
    );

    // Format the cost texts
    String monthlyText = defaultCurrency.formatAmount(monthlyCostInDefaultCurrency);
    String annualText = defaultCurrency.formatAmount(annualCostInDefaultCurrency);
    String totalSpentText = defaultCurrency.formatAmount(totalSpentInDefaultCurrency);

    // Add original amounts in parentheses if currencies differ
    if (subscriptionCurrency.code != defaultCurrency.code) {
      monthlyText += ' (${subscriptionCurrency.formatAmount(subscription.monthlyCost)})';
      annualText += ' (${subscriptionCurrency.formatAmount(subscription.annualCost)})';
      totalSpentText += ' (${subscriptionCurrency.formatAmount(subscription.totalAmountSpent)})';
    }

    return Card(
      margin: const EdgeInsets.only(bottom: 16),
      elevation: 2,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(16),
      ),
      child: InkWell(
        borderRadius: BorderRadius.circular(16),
        onTap: () {
          Navigator.push(
            context,
            MaterialPageRoute(
              builder: (context) =>
                  PaymentDetailScreen(subscription: subscription),
            ),
          );
        },
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Header with name and status badge
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Expanded(
                    child: Row(
                      children: [
                        // Subscription icon based on status
                        Container(
                          padding: const EdgeInsets.all(8),
                          decoration: BoxDecoration(
                            color: subscription.isActive 
                              ? Colors.green.withOpacity(0.1) 
                              : Colors.red.withOpacity(0.1),
                            borderRadius: BorderRadius.circular(12),
                          ),
                          child: Icon(
                            subscription.isActive 
                              ? Icons.check_circle_outline 
                              : Icons.cancel_outlined,
                            color: subscription.isActive ? Colors.green : Colors.red,
                            size: 20,
                          ),
                        ),
                        const SizedBox(width: 12),

                        // Subscription name
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                subscription.name,
                                style: const TextStyle(
                                  fontSize: 18,
                                  fontWeight: FontWeight.bold,
                                ),
                                overflow: TextOverflow.ellipsis,
                              ),
                              if (!subscription.isActive)
                                Text(
                                  "Cancelled",
                                  style: TextStyle(
                                    fontSize: 14,
                                    color: Colors.red.shade700,
                                    fontWeight: FontWeight.w500,
                                  ),
                                ),
                            ],
                          ),
                        ),
                        const SizedBox(width: 8),
                        _buildFrequencyBadge(subscription),
                      ],
                    ),
                  ),

                  // Menu button
                  PopupMenuButton<String>(
                    icon: Icon(
                      Icons.more_vert,
                      color: Theme.of(context).colorScheme.onSurface.withOpacity(0.7),
                    ),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                    onSelected: (value) async {
                      if (value == 'edit') {
                        Navigator.push(
                          context,
                          MaterialPageRoute(
                            builder: (context) => EditSubscriptionScreen(
                              subscription: subscription,
                            ),
                          ),
                        );
                      } else if (value == 'addPriceChange') {
                        showModalBottomSheet(
                          context: context,
                          isScrollControlled: true,
                          builder: (context) =>
                              PriceChangeForm(subscription: subscription),
                        );
                      } else if (value == 'delete') {
                        await DeleteSubscriptionDialog.show(
                          context: context,
                          subscription: subscription,
                        );
                      }
                    },
                    itemBuilder: (context) => [
                      PopupMenuItem<String>(
                        value: 'edit',
                        child: Row(
                          children: [
                            Icon(Icons.edit, color: Theme.of(context).colorScheme.primary),
                            const SizedBox(width: 8),
                            const Text('Edit'),
                          ],
                        ),
                      ),
                      PopupMenuItem<String>(
                        value: 'addPriceChange',
                        child: Row(
                          children: [
                            Icon(Icons.price_change, color: Colors.green.shade600),
                            const SizedBox(width: 8),
                            const Text('Add Price Change'),
                          ],
                        ),
                      ),
                      PopupMenuItem<String>(
                        value: 'delete',
                        child: Row(
                          children: [
                            Icon(Icons.delete, color: Colors.red.shade600),
                            const SizedBox(width: 8),
                            const Text('Delete'),
                          ],
                        ),
                      ),
                    ],
                  ),
                ],
              ),

              const Divider(height: 24),

              // Cost information
              if (subscription.isActive)
                Row(
                  children: [
                    // Monthly cost
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Row(
                            children: [
                              Icon(
                                Icons.calendar_today,
                                size: 14,
                                color: Theme.of(context).colorScheme.primary.withOpacity(0.7),
                              ),
                              const SizedBox(width: 4),
                              const Text(
                                'Monthly',
                                style: TextStyle(
                                  fontSize: 14,
                                  fontWeight: FontWeight.w500,
                                ),
                              ),
                            ],
                          ),
                          const SizedBox(height: 4),
                          Text(
                            monthlyText,
                            style: TextStyle(
                              fontSize: 16,
                              fontWeight: FontWeight.bold,
                              color: Theme.of(context).colorScheme.primary,
                            ),
                          ),
                        ],
                      ),
                    ),

                    // Annual cost
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Row(
                            children: [
                              Icon(
                                Icons.calendar_month,
                                size: 14,
                                color: Colors.purple.withOpacity(0.7),
                              ),
                              const SizedBox(width: 4),
                              const Text(
                                'Annually',
                                style: TextStyle(
                                  fontSize: 14,
                                  fontWeight: FontWeight.w500,
                                ),
                              ),
                            ],
                          ),
                          const SizedBox(height: 4),
                          Text(
                            annualText,
                            style: TextStyle(
                              fontSize: 16,
                              fontWeight: FontWeight.bold,
                              color: Colors.purple,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),

              // Next payment date
              if (subscription.isActive)
                Padding(
                  padding: const EdgeInsets.only(top: 16),
                  child: Row(
                    children: [
                      Icon(
                        Icons.event_available,
                        size: 16,
                        color: Theme.of(context).colorScheme.primary,
                      ),
                      const SizedBox(width: 8),
                      Text(
                        'Next payment: ${subscription.formattedNextPaymentDate}',
                        style: TextStyle(
                          fontSize: 14,
                          fontWeight: FontWeight.w500,
                          color: Theme.of(context).colorScheme.primary,
                        ),
                      ),
                    ],
                  ),
                ),

              // Total spent
              Padding(
                padding: const EdgeInsets.only(top: 8),
                child: Row(
                  children: [
                    Icon(
                      Icons.monetization_on,
                      size: 16,
                      color: Colors.green.shade600,
                    ),
                    const SizedBox(width: 8),
                    Text(
                      'Total spent: $totalSpentText',
                      style: TextStyle(
                        fontSize: 14,
                        fontWeight: FontWeight.w500,
                        color: Colors.green.shade600,
                      ),
                    ),
                  ],
                ),
              ),

              // Labels
              if (subscription.labels.isNotEmpty)
                Padding(
                  padding: const EdgeInsets.only(top: 12),
                  child: Wrap(
                    spacing: 6,
                    runSpacing: 6,
                    children: subscription.labels.map((label) {
                      final labelColor = Color(
                        int.parse(
                          label.color.substring(1, 7),
                          radix: 16,
                        ) + 0xFF000000,
                      );
                      return Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 8,
                          vertical: 4,
                        ),
                        decoration: BoxDecoration(
                          color: labelColor.withOpacity(0.2),
                          borderRadius: BorderRadius.circular(12),
                          border: Border.all(
                            color: labelColor,
                            width: 1,
                          ),
                        ),
                        child: Text(
                          label.name,
                          style: TextStyle(
                            fontSize: 12,
                            fontWeight: FontWeight.w500,
                            color: labelColor,
                          ),
                        ),
                      );
                    }).toList(),
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }
}
