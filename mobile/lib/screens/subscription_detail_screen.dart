import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:subscription_tracker/providers/family_member_provider.dart';
import 'package:subscription_tracker/providers/label_provider.dart';
import '../models/currency.dart';
import '../models/subscription.dart';
import '../models/subscription_state.dart';
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

        final defaultCurrency = subscriptionProvider.defaultCurrencyEnum;
        final subscriptionCurrency = Currency.fromCode(
          subscription.subscriptionPayments.isNotEmpty
              ? subscription.getLastPaymentDetail().currency
              : Currency.defaultCode,
        );
        final monthlyCostInDefaultCurrency = CurrencyConverter.convert(
          subscription.monthlyCost,
          subscriptionCurrency.code,
          defaultCurrency.code,
        );
        final annualCostInDefaultCurrency = monthlyCostInDefaultCurrency * 12;
        // Format the cost texts
        String monthlyText = defaultCurrency.formatAmount(
          monthlyCostInDefaultCurrency,
        );
        String annualText = defaultCurrency.formatAmount(
          annualCostInDefaultCurrency,
        );

        // Add original amounts in parentheses if currencies differ
        if (subscriptionCurrency.code != defaultCurrency.code) {
          monthlyText +=
              ' (${subscriptionCurrency.formatAmount(subscription.monthlyCost)})';
          annualText +=
              ' (${subscriptionCurrency.formatAmount(subscription.annualCost)})';
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
                  PopupMenuItem<String>(
                    value: 'edit',
                    child: Row(
                      children: [
                        Icon(
                          Icons.edit,
                          color: Theme.of(context).colorScheme.primary,
                        ),
                        SizedBox(width: 8),
                        Text('Edit'),
                      ],
                    ),
                  ),
                  if (subscription.state == SubscriptionState.active)
                    PopupMenuItem<String>(
                      value: 'addPaymentHistory',
                      child: Row(
                        children: [
                          Icon(
                            Icons.price_change,
                            color: Colors.green.shade600,
                          ),
                          SizedBox(width: 8),
                          Text('Price Change'),
                        ],
                      ),
                    ),
                  // Show stop subscription option if subscription is active
                  if (subscription.state == SubscriptionState.active)
                    PopupMenuItem<String>(
                      value: 'stopPayment',
                      child: Row(
                        children: [
                          Icon(
                            Icons.stop_circle,
                            color: Colors.orange.shade600,
                          ),
                          SizedBox(width: 8),
                          Text('Cancel'),
                        ],
                      ),
                    ),
                  // Show reactivate subscription option if subscription is not active
                  if (subscription.state == SubscriptionState.ended)
                    PopupMenuItem<String>(
                      value: 'reactivatePayment',
                      child: Row(
                        children: [
                          Icon(Icons.play_circle, color: Colors.green.shade600),
                          SizedBox(width: 8),
                          Text('Reactivate'),
                        ],
                      ),
                    ),
                  PopupMenuItem<String>(
                    value: 'delete',
                    child: Row(
                      children: [
                        Icon(Icons.delete_forever, color: Colors.red.shade600),
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
                Card(
                  elevation: 2,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(16),
                  ),
                  child: Padding(
                    padding: const EdgeInsets.all(16),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
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
                                        color: Theme.of(context)
                                            .colorScheme
                                            .primary
                                            .withValues(alpha: 0.7),
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
                                      color: Theme.of(
                                        context,
                                      ).colorScheme.primary,
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
                                        color: Colors.purple.withValues(
                                          alpha: 0.7,
                                        ),
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
                        const SizedBox(height: 12),
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
                                        Icons.attach_money,
                                        size: 14,
                                        color: Colors.green.withValues(
                                          alpha: 0.7,
                                        ),
                                      ),
                                      const SizedBox(width: 4),
                                      const Text(
                                        'Price',
                                        style: TextStyle(
                                          fontSize: 14,
                                          fontWeight: FontWeight.w500,
                                        ),
                                      ),
                                    ],
                                  ),
                                  const SizedBox(height: 4),
                                  Text(
                                    CurrencyConverter.formatAmountWithCurrency(
                                      currentDetail.price,
                                      currentDetail.currency,
                                    ),
                                    style: TextStyle(
                                      fontSize: 16,
                                      fontWeight: FontWeight.bold,
                                      color: Colors.green,
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
                                        color: Colors.blue.withValues(
                                          alpha: 0.7,
                                        ),
                                      ),
                                      const SizedBox(width: 4),
                                      const Text(
                                        'Type',
                                        style: TextStyle(
                                          fontSize: 14,
                                          fontWeight: FontWeight.w500,
                                        ),
                                      ),
                                    ],
                                  ),
                                  const SizedBox(height: 4),
                                  Text(
                                    subscriptionType,
                                    style: TextStyle(
                                      fontSize: 16,
                                      fontWeight: FontWeight.bold,
                                      color: Colors.blue,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 12),
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
                                        Icons.monetization_on,
                                        size: 14,
                                        color: Colors.lime.withValues(
                                          alpha: 0.7,
                                        ),
                                      ),
                                      const SizedBox(width: 4),
                                      const Text(
                                        'Total',
                                        style: TextStyle(
                                          fontSize: 14,
                                          fontWeight: FontWeight.w500,
                                        ),
                                      ),
                                    ],
                                  ),
                                  const SizedBox(height: 4),
                                  Text(
                                    subscription.formattedTotalAmountSpent,
                                    style: TextStyle(
                                      fontSize: 16,
                                      fontWeight: FontWeight.bold,
                                      color: Colors.lime,
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
                                        color: Colors.orange.withValues(
                                          alpha: 0.7,
                                        ),
                                      ),
                                      const SizedBox(width: 4),
                                      const Text(
                                        'Next',
                                        style: TextStyle(
                                          fontSize: 14,
                                          fontWeight: FontWeight.w500,
                                        ),
                                      ),
                                    ],
                                  ),
                                  const SizedBox(height: 4),
                                  Text(
                                    subscription.formattedNextPaymentDate,
                                    style: TextStyle(
                                      fontSize: 16,
                                      fontWeight: FontWeight.bold,
                                      color: Colors.orange,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),

                // Show labels if there are any
                if (subscription.labelIds.isNotEmpty)
                  Padding(
                    padding: const EdgeInsets.only(top: 12),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Padding(
                          padding: const EdgeInsets.only(left: 4, bottom: 12),
                          child: Text(
                            'Labels',
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.w600,
                              color: Theme.of(context).colorScheme.primary,
                            ),
                          ),
                        ),
                        Consumer<LabelProvider>(
                          builder: (context, provider, child) {
                            final labels = provider.labels
                                .where(
                                  (label) =>
                                      subscription.labelIds.contains(label.id),
                                )
                                .toList();
                            return Wrap(
                              spacing: 4,
                              runSpacing: 4,
                              children: labels.map((label) {
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
                            );
                          },
                        ),
                      ],
                    ),
                  ),

                // Show family members if any are assigned
                if (subscription.userFamilyMemberIds.isNotEmpty ||
                    subscription.payerFamilyMemberId != null)
                  Padding(
                    padding: const EdgeInsets.only(top: 12),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Padding(
                          padding: const EdgeInsets.only(left: 4, bottom: 12),
                          child: Text(
                            'Family',
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.w600,
                              color: Theme.of(context).colorScheme.primary,
                            ),
                          ),
                        ),
                        Card(
                          elevation: 2,
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(16),
                          ),
                          child: Padding(
                            padding: const EdgeInsets.all(16.0),
                            child: Column(
                              children: [
                                if (subscription.userFamilyMemberIds.isNotEmpty)
                                  Column(
                                    crossAxisAlignment:
                                        CrossAxisAlignment.start,
                                    children: [
                                      const Padding(
                                        padding: EdgeInsets.only(
                                          left: 16,
                                          top: 8,
                                          bottom: 8,
                                        ),
                                        child: Text(
                                          'Used by',
                                          style: TextStyle(
                                            fontSize: 16,
                                            fontWeight: FontWeight.w500,
                                          ),
                                        ),
                                      ),
                                      Consumer<FamilyMemberProvider>(
                                        builder: (context, provider, child) {
                                          final familyMembers = provider
                                              .familyMembers
                                              .where(
                                                (member) => subscription
                                                    .userFamilyMemberIds
                                                    .contains(member.id),
                                              )
                                              .toList();

                                          return Wrap(
                                            children: familyMembers
                                                .map(
                                                  (member) => ListTile(
                                                    leading: const CircleAvatar(
                                                      child: Icon(Icons.person),
                                                    ),
                                                    title: Text(
                                                      member.name,
                                                      style: TextStyle(
                                                        fontWeight:
                                                            FontWeight.bold,
                                                        color: Theme.of(
                                                          context,
                                                        ).colorScheme.primary,
                                                      ),
                                                    ),
                                                    dense: true,
                                                  ),
                                                )
                                                .toList(),
                                          );
                                        },
                                      ),
                                    ],
                                  ),
                                if (subscription.payerFamilyMemberId != null)
                                  Consumer<FamilyMemberProvider>(
                                    builder: (context, provider, child) {
                                      final payer = provider
                                          .getFamilyMemberById(
                                            subscription.payerFamilyMemberId,
                                          )!;
                                      return ListTile(
                                        leading: CircleAvatar(
                                          child: Icon(
                                            payer.id == 'family'
                                                ? Icons.group
                                                : Icons.account_balance_wallet,
                                          ),
                                        ),
                                        title: const Text('Paid by'),
                                        subtitle: Text(
                                          payer.name,
                                          style: TextStyle(
                                            fontWeight: FontWeight.bold,
                                            color: Theme.of(
                                              context,
                                            ).colorScheme.primary,
                                          ),
                                        ),
                                        dense: true,
                                      );
                                    },
                                  ),
                              ],
                            ),
                          ),
                        ),
                      ],
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
                          ).colorScheme.onSurface.withValues(alpha: 0.6),
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
                            history.state == SubscriptionState.notStarted
                                ? Icons.schedule
                                : history.state == SubscriptionState.active
                                ? Icons.check_circle
                                : Icons.stop_circle,
                            color: history.state == SubscriptionState.notStarted
                                ? Colors.orange
                                : history.state == SubscriptionState.active
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
                            history.state == SubscriptionState.notStarted
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
