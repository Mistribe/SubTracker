import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../models/payment.dart';
import '../providers/payment_provider.dart';
import '../widgets/edit_payment_form.dart';

class PaymentDetailScreen extends StatefulWidget {
  final Payment payment;

  const PaymentDetailScreen({super.key, required this.payment});

  @override
  State<PaymentDetailScreen> createState() => _PaymentDetailScreenState();
}

class _PaymentDetailScreenState extends State<PaymentDetailScreen> {
  late Payment payment;

  @override
  void initState() {
    super.initState();
    payment = widget.payment;
  }

  // Show dialog to stop a payment
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
                  'Are you sure you want to stop this payment? '
                  'A stopped payment will not be counted in active or total payments.',
                ),
                const SizedBox(height: 16),
                SwitchListTile(
                  title: const Text('Use last payment date'),
                  subtitle: const Text(
                    'If enabled, the payment will be stopped at the last payment date',
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
                        firstDate: payment.getLastPaymentDetail().startDate,
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
                        content: Text('Stopping payment...'),
                        duration: Duration(seconds: 1),
                      ),
                    );

                    // Stop the payment using the provider
                    await Provider.of<PaymentProvider>(
                      context,
                      listen: false,
                    ).stopPayment(
                      payment.id,
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

  // Show dialog to reactivate a payment
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
                  'Select a date to reactivate this payment. '
                  'The payment will be reactivated on this date.',
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

                    // Reactivate the payment using the provider
                    await Provider.of<PaymentProvider>(
                      context,
                      listen: false,
                    ).reactivatePayment(payment.id, selectedDate);

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

  // Show dialog to add a payment history entry
  void _showAddPaymentHistoryDialog(BuildContext context) {
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
                if (payment.isAnnual)
                  // For yearly payments, only allow changing the year
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
                  // For monthly payments, allow changing year and month
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
                        firstDate: payment.paymentDate,
                        lastDate: DateTime(2101),
                        selectableDayPredicate: (DateTime date) {
                          // Only allow selecting the same day of the month as the payment date
                          return date.day == payment.paymentDate.day;
                        },
                      );
                      if (picked != null && picked != selectedDate) {
                        setState(() {
                          // Ensure we keep the same day as the payment date
                          selectedDate = DateTime(
                            picked.year,
                            picked.month,
                            payment.paymentDate.day,
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

                      // Add the payment history entry using the provider
                      await Provider.of<PaymentProvider>(
                        context,
                        listen: false,
                      ).addPaymentDetailEntry(
                        payment.id,
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
    return Consumer<PaymentProvider>(
      builder: (context, paymentProvider, child) {
        // Get the updated payment from the provider
        final updatedPayment = paymentProvider.payments.firstWhere(
          (p) => p.id == payment.id,
          orElse: () => payment,
        );

        // Update the local payment if it has changed
        if (updatedPayment != payment) {
          payment = updatedPayment;
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
                  // Show stop payment option if payment is not stopped
                  if (!payment.isStopped)
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
                  // Show reactivate payment option if payment is stopped
                  if (payment.isStopped)
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
                          'Payment Summary',
                          style: Theme.of(context).textTheme.titleLarge,
                        ),
                        const SizedBox(height: 16),
                        buildInfoRow(
                          'Current Price:',
                          '\$${payment.price.toStringAsFixed(2)}',
                          Icons.attach_money,
                        ),
                        buildInfoRow(
                          'Payment Type:',
                          payment.isAnnual ? 'Annual' : 'Monthly',
                          Icons.calendar_today,
                        ),
                        if (payment.isAnnual)
                          buildInfoRow(
                            'Monthly Cost:',
                            '\$${payment.monthlyCost.toStringAsFixed(2)}',
                            Icons.calculate,
                          ),
                        buildInfoRow(
                          'Start Date:',
                          '${payment.paymentDate.month}/${payment.paymentDate.day}/${payment.paymentDate.year}',
                          Icons.date_range,
                        ),
                        buildInfoRow(
                          'Next Payment:',
                          payment.formattedNextPaymentDate,
                          Icons.event_available,
                        ),
                        buildInfoRow(
                          'Total Spent:',
                          payment.formattedTotalAmountSpent,
                          Icons.monetization_on,
                        ),
                        // Show payment status (active or stopped)
                        buildInfoRow(
                          'Status:',
                          payment.isStopped ? 'Stopped' : 'Active',
                          payment.isStopped
                              ? Icons.stop_circle
                              : Icons.check_circle,
                        ),
                        // Show stop date if payment is stopped and has a stop date
                        if (payment.isStopped && payment.stopDate != null)
                          buildInfoRow(
                            'Stop Date:',
                            '${payment.stopDate!.month}/${payment.stopDate!.day}/${payment.stopDate!.year}',
                            Icons.stop_circle,
                          ),
                        // Show reactivation date if payment is stopped and has a reactivation date
                        if (payment.isStopped &&
                            payment.reactivationDate != null)
                          buildInfoRow(
                            'Reactivation Date:',
                            '${payment.reactivationDate!.month}/${payment.reactivationDate!.day}/${payment.reactivationDate!.year}',
                            Icons.event_available,
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
                    ElevatedButton.icon(
                      icon: const Icon(Icons.add),
                      label: const Text('Add Payment History'),
                      onPressed: () {
                        _showAddPaymentHistoryDialog(context);
                      },
                    ),
                  ],
                ),
                const SizedBox(height: 16),

                if (payment.paymentDetails.isEmpty)
                  const Card(
                    child: Padding(
                      padding: EdgeInsets.all(16.0),
                      child: Text(
                        'No payment history recorded yet.',
                        style: TextStyle(fontSize: 16),
                      ),
                    ),
                  )
                else
                  Card(
                    child: ListView.separated(
                      shrinkWrap: true,
                      physics: const NeverScrollableScrollPhysics(),
                      itemCount: payment.paymentDetails.length,
                      separatorBuilder: (context, index) => const Divider(),
                      itemBuilder: (context, index) {
                        final history = payment.paymentDetails[index];
                        return ListTile(
                          leading: Icon(
                            history.isActive
                                ? Icons.check_circle
                                : Icons.stop_circle,
                            color: history.isActive ? Colors.green : Colors.red,
                          ),
                          title: Text('\$${history.price.toStringAsFixed(2)}'),
                          subtitle: Text(
                            'from ${history.startDate.month}/${history.startDate.day}/${history.startDate.year} to ${history.endDate.year == 9999 ? "now" : "${history.endDate.month}/${history.endDate.day}/${history.endDate.year}"}',
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
