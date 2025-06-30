import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../models/payment.dart';
import '../providers/payment_provider.dart';
import '../widgets/edit_payment_form.dart';

class PaymentDetailScreen extends StatelessWidget {
  final Payment payment;

  const PaymentDetailScreen({super.key, required this.payment});

  // Show dialog to edit a price change
  void _showEditPriceChangeDialog(BuildContext context, int priceChangeIndex) {
    final priceChange = payment.priceHistory[priceChangeIndex];
    final priceController = TextEditingController(text: priceChange.price.toString());
    DateTime selectedDate = priceChange.endDate;

    showDialog(
      context: context,
      builder: (ctx) => StatefulBuilder(
        builder: (context, setState) {
          return AlertDialog(
            title: const Text('Edit Price Change'),
            content: Column(
              mainAxisSize: MainAxisSize.min,
              children: [
                TextFormField(
                  controller: priceController,
                  decoration: const InputDecoration(
                    labelText: 'Price',
                    hintText: 'Enter the price',
                    prefixIcon: Icon(Icons.attach_money),
                  ),
                  keyboardType: const TextInputType.numberWithOptions(decimal: true),
                  inputFormatters: [
                    FilteringTextInputFormatter.allow(RegExp(r'^\d+\.?\d{0,2}')),
                  ],
                ),
                const SizedBox(height: 16),
                ListTile(
                  leading: const Icon(Icons.calendar_today),
                  title: const Text('Effective Date'),
                  subtitle: Text('${selectedDate.month}/${selectedDate.day}/${selectedDate.year}'),
                  onTap: () async {
                    final DateTime? picked = await showDatePicker(
                      context: context,
                      initialDate: selectedDate,
                      firstDate: payment.paymentDate,
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
                          const SnackBar(content: Text('Price must be greater than zero')),
                        );
                        return;
                      }

                      // Update the price change using the provider
                      Provider.of<PaymentProvider>(context, listen: false)
                          .updatePriceChange(payment.id, priceChangeIndex, newPrice, selectedDate);

                      Navigator.of(context).pop();

                      // Show success message
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(
                          content: Text('Price change updated successfully'),
                          duration: Duration(seconds: 2),
                        ),
                      );
                    } catch (e) {
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(content: Text('Please enter a valid number')),
                      );
                    }
                  } else {
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(content: Text('Please enter a price')),
                    );
                  }
                },
                child: const Text('Update'),
              ),
            ],
          );
        },
      ),
    );
  }

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
                  keyboardType: const TextInputType.numberWithOptions(decimal: true),
                  inputFormatters: [
                    FilteringTextInputFormatter.allow(RegExp(r'^\d+\.?\d{0,2}')),
                  ],
                ),
                const SizedBox(height: 16),
                ListTile(
                  leading: const Icon(Icons.calendar_today),
                  title: const Text('Effective Date'),
                  subtitle: Text('${selectedDate.month}/${selectedDate.day}/${selectedDate.year}'),
                  onTap: () async {
                    final DateTime? picked = await showDatePicker(
                      context: context,
                      initialDate: selectedDate,
                      firstDate: payment.paymentDate,
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
                          const SnackBar(content: Text('Price must be greater than zero')),
                        );
                        return;
                      }

                      // Add the price change using the provider
                      Provider.of<PaymentProvider>(context, listen: false)
                          .addPriceChange(payment.id, newPrice, selectedDate);

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
                        const SnackBar(content: Text('Please enter a valid number')),
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
    return Scaffold(
      appBar: AppBar(
        title: Text(payment.name),
        actions: [
          PopupMenuButton<String>(
            icon: const Icon(Icons.more_vert),
            onSelected: (value) {
              if (value == 'edit') {
                showModalBottomSheet(
                  context: context,
                  isScrollControlled: true,
                  builder: (context) => EditPaymentForm(payment: payment),
                );
              } else if (value == 'addPriceChange') {
                _showAddPriceChangeDialog(context);
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
                value: 'addPriceChange',
                child: Row(
                  children: [
                    Icon(Icons.price_change),
                    SizedBox(width: 8),
                    Text('Add Price Change'),
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
                      'Payment Summary',
                      style: Theme.of(context).textTheme.titleLarge,
                    ),
                    const SizedBox(height: 16),
                    _buildInfoRow(
                      context,
                      'Current Price:',
                      '\$${payment.price.toStringAsFixed(2)}',
                      Icons.attach_money,
                    ),
                    _buildInfoRow(
                      context,
                      'Payment Type:',
                      payment.isAnnual ? 'Annual' : 'Monthly',
                      Icons.calendar_today,
                    ),
                    if (payment.isAnnual)
                      _buildInfoRow(
                        context,
                        'Monthly Cost:',
                        '\$${payment.monthlyCost.toStringAsFixed(2)}',
                        Icons.calculate,
                      ),
                    _buildInfoRow(
                      context,
                      'Start Date:',
                      '${payment.paymentDate.month}/${payment.paymentDate.day}/${payment.paymentDate.year}',
                      Icons.date_range,
                    ),
                    _buildInfoRow(
                      context,
                      'Next Payment:',
                      payment.formattedNextPaymentDate,
                      Icons.event_available,
                    ),
                    _buildInfoRow(
                      context,
                      'Total Spent:',
                      payment.formattedTotalAmountSpent,
                      Icons.monetization_on,
                    ),
                  ],
                ),
              ),
            ),

            const SizedBox(height: 24),

            // Price history section
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  'Price History',
                  style: Theme.of(context).textTheme.titleLarge,
                ),
                ElevatedButton.icon(
                  icon: const Icon(Icons.add),
                  label: const Text('Add Price Change'),
                  onPressed: () {
                    _showAddPriceChangeDialog(context);
                  },
                ),
              ],
            ),
            const SizedBox(height: 16),

            if (payment.priceHistory.isEmpty)
              const Card(
                child: Padding(
                  padding: EdgeInsets.all(16.0),
                  child: Text(
                    'No price changes recorded yet.',
                    style: TextStyle(fontSize: 16),
                  ),
                ),
              )
            else
              Card(
                child: ListView.separated(
                  shrinkWrap: true,
                  physics: const NeverScrollableScrollPhysics(),
                  itemCount: payment.priceHistory.length,
                  separatorBuilder: (context, index) => const Divider(),
                  itemBuilder: (context, index) {
                    final previousChangeDate = index == payment.priceHistory.length - 1
                    ? payment.paymentDate : payment.priceHistory[index + 1].endDate;
                    final priceChange = payment.priceHistory[index];
                    return ListTile(
                      leading: const Icon(Icons.history),
                      title: Text('\$${priceChange.price.toStringAsFixed(2)}'),
                      subtitle: Text(
                        'from ${previousChangeDate.month}/${previousChangeDate.day}/${previousChangeDate.year} to ${priceChange.endDate.month}/${priceChange.endDate.day}/${priceChange.endDate.year}',
                      ),
                      trailing: IconButton(
                        icon: const Icon(Icons.edit),
                        onPressed: () {
                          _showEditPriceChangeDialog(context, index);
                        },
                        tooltip: 'Edit Price Change',
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
  }

  Widget _buildInfoRow(BuildContext context, String label, String value, IconData icon) {
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
}
