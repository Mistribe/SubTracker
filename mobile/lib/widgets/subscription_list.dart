import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../providers/subscription_provider.dart';
import '../models/subscription.dart';
import '../screens/subscription_detail_screen.dart';
import 'edit_payment_form.dart';

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

  // Show dialog to add a price change
  void _showAddPriceChangeDialog(BuildContext context) {
    final priceController = TextEditingController();
    DateTime selectedDate = DateTime.now();

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
                ListTile(
                  leading: const Icon(Icons.calendar_today),
                  title: const Text('Effective Date'),
                  subtitle: Text(
                    '${selectedDate.month}/${selectedDate.day}/${selectedDate.year}',
                  ),
                  onTap: () async {
                    final DateTime? picked = await showDatePicker(
                      context: context,
                      initialDate: selectedDate,
                      firstDate: subscription.getLastPaymentDetail().startDate,
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
                onPressed: () {
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

                      // Add the price change using the provider
                      Provider.of<SubscriptionProvider>(
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
                        const SnackBar(
                          content: Text('Please enter a valid number'),
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
    final subscriptionProvider = Provider.of<SubscriptionProvider>(
      context,
      listen: false,
    );

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
                      Text(
                        'Monthly: \$${subscription.monthlyCost.toStringAsFixed(2)}',
                        style: const TextStyle(
                          fontSize: 14,
                          color: Colors.grey,
                        ),
                      ),
                      Text(
                        'Annually: \$${subscription.annualCost.toStringAsFixed(2)}',
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
                            'Total spent: ${subscription.formattedTotalAmountSpent}',
                            style: const TextStyle(
                              fontSize: 14,
                              fontWeight: FontWeight.w500,
                              color: Colors.green,
                            ),
                          ),
                        ],
                      ),
                    ],
                  ),
                  PopupMenuButton<String>(
                    icon: const Icon(Icons.more_vert),
                    onSelected: (value) async {
                      if (value == 'edit') {
                        showModalBottomSheet(
                          context: context,
                          isScrollControlled: true,
                          builder: (context) =>
                              EditSubscriptionForm(subscription: subscription),
                        );
                      } else if (value == 'addPriceChange') {
                        _showAddPriceChangeDialog(context);
                      } else if (value == 'delete') {
                        try {
                          // Show loading indicator
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(
                              content: Text('Removing subscription...'),
                              duration: Duration(seconds: 1),
                            ),
                          );

                          // Remove the subscription
                          await subscriptionProvider.removePayment(
                            subscription.id,
                          );

                          // Show success message
                          ScaffoldMessenger.of(context).showSnackBar(
                            SnackBar(
                              content: Text('${subscription.name} removed'),
                              duration: const Duration(seconds: 2),
                            ),
                          );
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
                        }
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
