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
    final subscriptions = subscriptionProvider.subscriptions;

    if (subscriptions.isEmpty) {
      return const Center(
        child: Text(
          'No subscriptions yet. Add your first subscription!',
          style: TextStyle(fontSize: 16),
        ),
      );
    }

    return ListView.builder(
      padding: const EdgeInsets.all(16),
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

  @override
  Widget build(BuildContext context) {
    final subscriptionProvider = Provider.of<SubscriptionProvider>(
      context,
      listen: false,
    );

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
    String monthlyText = 'Monthly: ${defaultCurrency.formatAmount(monthlyCostInDefaultCurrency)}';
    String annualText = 'Annually: ${defaultCurrency.formatAmount(annualCostInDefaultCurrency)}';
    String totalSpentText = 'Total spent: ${defaultCurrency.formatAmount(totalSpentInDefaultCurrency)}';

    // Add original amounts in parentheses if currencies differ
    if (subscriptionCurrency.code != defaultCurrency.code) {
      monthlyText += ' (${subscriptionCurrency.formatAmount(subscription.monthlyCost)})';
      annualText += ' (${subscriptionCurrency.formatAmount(subscription.annualCost)})';
      totalSpentText += ' (${subscriptionCurrency.formatAmount(subscription.totalAmountSpent)})';
    }

    return Card(
      margin: const EdgeInsets.only(bottom: 16),
      child: InkWell(
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
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Expanded(
                    child: Text(
                      subscription.name,
                      style: const TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                      ),
                      overflow: TextOverflow.ellipsis,
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 8),
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      if (subscription.isActive)
                        Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              monthlyText,
                              style: const TextStyle(
                                fontSize: 14,
                                color: Colors.grey,
                              ),
                            ),
                            Text(
                              annualText,
                              style: const TextStyle(fontSize: 16),
                            ),
                            const SizedBox(height: 4),
                            Row(
                              children: [
                                const Icon(
                                  Icons.event_available,
                                  size: 14,
                                  color: Colors.blue,
                                ),
                                const SizedBox(width: 4),
                                Text(
                                  'Next payment: ${subscription.formattedNextPaymentDate}',
                                  style: const TextStyle(
                                    fontSize: 14,
                                    fontWeight: FontWeight.w500,
                                    color: Colors.blue,
                                  ),
                                ),
                              ],
                            ),
                          ],
                        ),
                      if (!subscription.isActive) Text("Cancelled"),
                      const SizedBox(height: 4),
                      Row(
                        children: [
                          const Icon(
                            Icons.monetization_on,
                            size: 14,
                            color: Colors.green,
                          ),
                          const SizedBox(width: 4),
                          Text(
                            totalSpentText,
                            style: const TextStyle(
                              fontSize: 14,
                              fontWeight: FontWeight.w500,
                              color: Colors.green,
                            ),
                          ),
                        ],
                      ),
                      // Display labels if there are any
                      if (subscription.labels.isNotEmpty)
                        Padding(
                          padding: const EdgeInsets.only(top: 8.0),
                          child: Wrap(
                            spacing: 4,
                            runSpacing: 4,
                            children: subscription.labels.map((label) {
                              return Chip(
                                label: Text(
                                  label.name,
                                  style: const TextStyle(
                                    fontSize: 10,
                                    color: Colors.white,
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
                                  horizontal: 2,
                                ),
                                labelPadding: const EdgeInsets.symmetric(
                                  horizontal: 2,
                                ),
                              );
                            }).toList(),
                          ),
                        ),
                    ],
                  ),
                  PopupMenuButton<String>(
                    icon: const Icon(Icons.more_vert),
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
                      const PopupMenuItem<String>(
                        value: 'edit',
                        child: Row(
                          children: [
                            Icon(Icons.edit, color: Colors.blue),
                            SizedBox(width: 8),
                            Text('Edit'),
                          ],
                        ),
                      ),
                      const PopupMenuItem<String>(
                        value: 'addPriceChange',
                        child: Row(
                          children: [
                            Icon(Icons.price_change, color: Colors.green),
                            SizedBox(width: 8),
                            Text('Add Price Change'),
                          ],
                        ),
                      ),
                      const PopupMenuItem<String>(
                        value: 'delete',
                        child: Row(
                          children: [
                            Icon(Icons.delete, color: Colors.red),
                            SizedBox(width: 8),
                            Text('Delete'),
                          ],
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}
