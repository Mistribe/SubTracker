import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../models/subscription.dart';
import '../providers/subscription_provider.dart';
import '../widgets/edit_payment_form.dart';

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

  // Show dialog to stop a subscription
  void _showStopPaymentDialog(BuildContext context) {
    // Default stop date is the current date
    DateTime selectedDate = DateTime.now();
    bool useLastPaymentDate = true;

    showDialog(
      context: context,
      builder: (ctx) => StatefulBuilder(
        builder: (context, setState) {
          return AlertDialog(
            title: const Text('Stop Payment'),
            content: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
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
                  value: useLastPaymentDate,
                  onChanged: (value) {
                    setState(() {
                      useLastPaymentDate = value;
                    });
                  },
                ),
                if (!useLastPaymentDate)
                  ListTile(
                    leading: const Icon(Icons.calendar_today),
                    title: const Text('Stop Date'),
                    subtitle: Text(
                      '${selectedDate.month}/${selectedDate.day}/${selectedDate.year}',
                    ),
                    onTap: () async {
                      final DateTime? picked = await showDatePicker(
                        context: context,
                        initialDate: selectedDate,
                        firstDate: subscription
                            .getLastPaymentDetail()
                            .startDate,
                        lastDate: DateTime.now(),
                      );
                      if (picked != null && picked != selectedDate) {
                        setState(() {
                          selectedDate = picked;
                        });
                      }
                    },
                  ),
              ],
            ),
            actions: [
              TextButton(
                onPressed: () {
                  Navigator.of(context).pop();
                },
                child: const Text('Cancel'),
              ),
              TextButton(
                onPressed: () async {
                  try {
                    // Show loading indicator
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(
                        content: Text('Stopping subscription...'),
                        duration: Duration(seconds: 1),
                      ),
                    );

                    // Stop the subscription using the provider
                    await Provider.of<SubscriptionProvider>(
                      context,
                      listen: false,
                    ).stopPayment(
                      subscription.id,
                      stopDate: useLastPaymentDate ? null : selectedDate,
                    );

                    Navigator.of(context).pop();

                    // Show success message
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(
                        content: Text('Payment stopped successfully'),
                        duration: Duration(seconds: 2),
                      ),
                    );
                  } catch (e) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(
                        content: Text('Error: ${e.toString()}'),
                        backgroundColor: Colors.red,
                        duration: const Duration(seconds: 3),
                      ),
                    );
                  }
                },
                child: const Text('Stop Payment'),
              ),
            ],
          );
        },
      ),
    );
  }

  // Show dialog to reactivate a subscription
  void _showReactivatePaymentDialog(BuildContext context) {
    DateTime selectedDate = DateTime.now();

    showDialog(
      context: context,
      builder: (ctx) => StatefulBuilder(
        builder: (context, setState) {
          return AlertDialog(
            title: const Text('Reactivate Payment'),
            content: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                const Text(
                  'Select a date to reactivate this subscription. '
                  'The subscription will be reactivated on this date.',
                ),
                const SizedBox(height: 16),
                ListTile(
                  leading: const Icon(Icons.calendar_today),
                  title: const Text('Reactivation Date'),
                  subtitle: Text(
                    '${selectedDate.month}/${selectedDate.day}/${selectedDate.year}',
                  ),
                  onTap: () async {
                    final DateTime? picked = await showDatePicker(
                      context: context,
                      initialDate: selectedDate,
                      firstDate: DateTime.now(),
                      lastDate: DateTime(2101),
                    );
                    if (picked != null && picked != selectedDate) {
                      setState(() {
                        selectedDate = picked;
                      });
                    }
                  },
                ),
              ],
            ),
            actions: [
              TextButton(
                onPressed: () {
                  Navigator.of(context).pop();
                },
                child: const Text('Cancel'),
              ),
              TextButton(
                onPressed: () async {
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
                    ).reactivatePayment(subscription.id, selectedDate);

                    Navigator.of(context).pop();

                    // Show success message
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(
                        content: Text(
                          'Payment reactivation scheduled successfully',
                        ),
                        duration: Duration(seconds: 2),
                      ),
                    );
                  } catch (e) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(
                        content: Text('Error: ${e.toString()}'),
                        backgroundColor: Colors.red,
                        duration: const Duration(seconds: 3),
                      ),
                    );
                  }
                },
                child: const Text('Reactivate'),
              ),
            ],
          );
        },
      ),
    );
  }

  // Show dialog to add a subscription history entry
  void _showAddPaymentHistoryDialog(BuildContext context) {
    final priceController = TextEditingController();
    DateTime selectedDate = DateTime.now();
    final currentDetail = subscription.getLastPaymentDetail();

    showDialog(
      context: context,
      builder: (ctx) => StatefulBuilder(
        builder: (context, setState) {
          return AlertDialog(
            title: const Text('Add Price Change'),
            content: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                TextFormField(
                  controller: priceController,
                  decoration: const InputDecoration(
                    labelText: 'New Price',
                    hintText: 'Enter the new price',
                    prefixIcon: Icon(Icons.attach_money),
                  ),
                  keyboardType: const TextInputType.numberWithOptions(
                    decimal: true,
                  ),
                  inputFormatters: [
                    FilteringTextInputFormatter.allow(
                      RegExp(r'^\d+\.?\d{0,2}'),
                    ),
                  ],
                ),
                const SizedBox(height: 16),
                if (currentDetail.months == 12)
                  // For yearly subscriptions, only allow changing the year
                  ListTile(
                    leading: const Icon(Icons.calendar_today),
                    title: const Text('Effective Year'),
                    subtitle: Text('${selectedDate.year}'),
                    onTap: () async {
                      // Show a dialog to select only the year
                      final int? selectedYear = await showDialog<int>(
                        context: context,
                        builder: (BuildContext context) {
                          int year = selectedDate.year;
                          return AlertDialog(
                            title: const Text('Select Year'),
                            content: SizedBox(
                              height: 300,
                              width: 300,
                              child: ListView.builder(
                                itemCount: 30,
                                // Show 30 years from current year
                                itemBuilder: (context, index) {
                                  final int yearOption =
                                      DateTime.now().year + index;
                                  return ListTile(
                                    title: Text(yearOption.toString()),
                                    onTap: () {
                                      Navigator.of(context).pop(yearOption);
                                    },
                                  );
                                },
                              ),
                            ),
                          );
                        },
                      );

                      if (selectedYear != null) {
                        setState(() {
                          // Keep the same month and day, only change the year
                          selectedDate = DateTime(
                            selectedYear,
                            selectedDate.month,
                            selectedDate.day,
                          );
                        });
                      }
                    },
                  )
                else
                  // For non-yearly subscriptions, allow changing year and month
                  ListTile(
                    leading: const Icon(Icons.calendar_today),
                    title: const Text('Effective Month/Year'),
                    subtitle: Text(
                      '${selectedDate.month}/${selectedDate.year}',
                    ),
                    onTap: () async {
                      // Show a dialog to select month and year
                      final DateTime? picked = await showDatePicker(
                        context: context,
                        initialDate: selectedDate,
                        firstDate: currentDetail.startDate,
                        lastDate: DateTime(2101),
                        selectableDayPredicate: (DateTime date) {
                          // Only allow selecting the same day of the month as the subscription date
                          return date.day == currentDetail.startDate.day;
                        },
                      );
                      if (picked != null && picked != selectedDate) {
                        setState(() {
                          // Ensure we keep the same day as the subscription date
                          selectedDate = DateTime(
                            picked.year,
                            picked.month,
                            currentDetail.startDate.day,
                          );
                        });
                      }
                    },
                  ),
              ],
            ),
            actions: [
              TextButton(
                onPressed: () {
                  Navigator.of(context).pop();
                },
                child: const Text('Cancel'),
              ),
              TextButton(
                onPressed: () async {
                  if (priceController.text.isNotEmpty) {
                    try {
                      final newPrice = double.parse(priceController.text);
                      if (newPrice <= 0) {
                        ScaffoldMessenger.of(context).showSnackBar(
                          const SnackBar(
                            content: Text('Price must be greater than zero'),
                          ),
                        );
                        return;
                      }

                      // Show loading indicator
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(
                          content: Text('Adding price change...'),
                          duration: Duration(seconds: 1),
                        ),
                      );

                      // Add the subscription history entry using the provider
                      await Provider.of<SubscriptionProvider>(
                        context,
                        listen: false,
                      ).addPaymentDetailEntry(
                        subscription.id,
                        newPrice,
                        selectedDate,
                      );

                      Navigator.of(context).pop();

                      // Show success message
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(
                          content: Text('Price change added successfully'),
                          duration: Duration(seconds: 2),
                        ),
                      );
                    } catch (e) {
                      ScaffoldMessenger.of(context).showSnackBar(
                        SnackBar(
                          content: Text('Error: ${e.toString()}'),
                          backgroundColor: Colors.red,
                          duration: const Duration(seconds: 3),
                        ),
                      );
                    }
                  } else {
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(content: Text('Please enter a price')),
                    );
                  }
                },
                child: const Text('Add'),
              ),
            ],
          );
        },
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
          return Padding(
            padding: const EdgeInsets.only(bottom: 12.0),
            child: Row(
              children: [
                Icon(icon, size: 20, color: Theme.of(context).primaryColor),
                const SizedBox(width: 8),
                Text(
                  label,
                  style: const TextStyle(
                    fontWeight: FontWeight.bold,
                    fontSize: 16,
                  ),
                ),
                const SizedBox(width: 8),
                Expanded(
                  child: Text(
                    value,
                    style: const TextStyle(fontSize: 16),
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
            title: Text(subscription.name),
            actions: [
              PopupMenuButton<String>(
                icon: const Icon(Icons.more_vert),
                onSelected: (value) {
                  if (value == 'edit') {
                    showModalBottomSheet(
                      context: context,
                      isScrollControlled: true,
                      builder: (context) =>
                          EditSubscriptionForm(subscription: subscription),
                    );
                  } else if (value == 'addPaymentHistory') {
                    _showAddPaymentHistoryDialog(context);
                  } else if (value == 'stopPayment') {
                    _showStopPaymentDialog(context);
                  } else if (value == 'reactivatePayment') {
                    _showReactivatePaymentDialog(context);
                  }
                },
                itemBuilder: (BuildContext context) => <PopupMenuEntry<String>>[
                  const PopupMenuItem<String>(
                    value: 'edit',
                    child: Row(
                      children: [
                        Icon(Icons.edit),
                        SizedBox(width: 8),
                        Text('Edit Payment'),
                      ],
                    ),
                  ),
                  const PopupMenuItem<String>(
                    value: 'addPaymentHistory',
                    child: Row(
                      children: [
                        Icon(Icons.price_change),
                        SizedBox(width: 8),
                        Text('Add Price Change'),
                      ],
                    ),
                  ),
                  // Show stop subscription option if subscription is active
                  if (subscription.isActive)
                    const PopupMenuItem<String>(
                      value: 'stopPayment',
                      child: Row(
                        children: [
                          Icon(Icons.stop_circle),
                          SizedBox(width: 8),
                          Text('Stop Payment'),
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
                          Text('Reactivate Payment'),
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
                // Payment summary card
                Card(
                  child: Padding(
                    padding: const EdgeInsets.all(16.0),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          'Subscription Summary',
                          style: Theme.of(context).textTheme.titleLarge,
                        ),
                        const SizedBox(height: 16),
                        buildInfoRow(
                          'Current Price:',
                          '\$${currentDetail.price.toStringAsFixed(2)}',
                          Icons.attach_money,
                        ),
                        buildInfoRow(
                          'Payment Type:',
                          subscriptionType,
                          Icons.calendar_today,
                        ),
                        // Always show monthly cost for non-monthly subscriptions
                        if (currentDetail.months > 1)
                          buildInfoRow(
                            'Monthly Cost:',
                            '\$${subscription.monthlyCost.toStringAsFixed(2)}',
                            Icons.calculate,
                          ),
                        buildInfoRow(
                          'Start Date:',
                          '${currentDetail.startDate.month}/${currentDetail.startDate.day}/${currentDetail.startDate.year}',
                          Icons.date_range,
                        ),
                        buildInfoRow(
                          'Next Payment:',
                          subscription.formattedNextPaymentDate,
                          Icons.event_available,
                        ),
                        buildInfoRow(
                          'Total Spent:',
                          subscription.formattedTotalAmountSpent,
                          Icons.monetization_on,
                        ),
                        // Show subscription status (active or stopped)
                        buildInfoRow(
                          'Status:',
                          subscription.isActive ? 'Active' : 'Stopped',
                          subscription.isActive
                              ? Icons.check_circle
                              : Icons.stop_circle,
                        ),
                        // Show end date if subscription is not active and has an end date
                        if (!subscription.isActive &&
                            currentDetail.endDate != null)
                          buildInfoRow(
                            'Stop Date:',
                            '${currentDetail.endDate!.month}/${currentDetail.endDate!.day}/${currentDetail.endDate!.year}',
                            Icons.stop_circle,
                          ),
                      ],
                    ),
                  ),
                ),

                const SizedBox(height: 24),

                // Payment History section
                Row(
                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                  children: [
                    Text(
                      'Payment History',
                      style: Theme.of(context).textTheme.titleLarge,
                    ),
                  ],
                ),
                const SizedBox(height: 16),

                if (subscription.subscriptionPayments.isEmpty)
                  const Card(
                    child: Padding(
                      padding: EdgeInsets.all(16.0),
                      child: Text(
                        'No subscription history recorded yet.',
                        style: TextStyle(fontSize: 16),
                      ),
                    ),
                  )
                else
                  Card(
                    child: ListView.separated(
                      shrinkWrap: true,
                      physics: const NeverScrollableScrollPhysics(),
                      itemCount: subscription.subscriptionPayments.length,
                      separatorBuilder: (context, index) => const Divider(),
                      itemBuilder: (context, index) {
                        final history =
                            subscription.subscriptionPayments[index];
                        return ListTile(
                          leading: Icon(
                            history.isActive
                                ? Icons.check_circle
                                : Icons.stop_circle,
                            color: history.isActive ? Colors.green : Colors.red,
                          ),
                          title: Text('\$${history.price.toStringAsFixed(2)}'),
                          subtitle: Text(
                            'from ${history.startDate.month}/${history.startDate.day}/${history.startDate.year} to ${history.endDate == null ? "now" : "${history.endDate!.month}/${history.endDate!.day}/${history.endDate!.year}"}',
                          ),
                          trailing: Text(
                            history.isActive ? 'Active' : 'Stopped',
                            style: TextStyle(
                              color: history.isActive
                                  ? Colors.green
                                  : Colors.red,
                              fontWeight: FontWeight.bold,
                            ),
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
