import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/payment_provider.dart';
import '../providers/theme_provider.dart';
import '../widgets/payment_list.dart';
import '../widgets/add_payment_form.dart';
import 'settings_screen.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final paymentProvider = Provider.of<PaymentProvider>(context);
    final themeProvider = Provider.of<ThemeProvider>(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Recurrent Payment Tracker'),
        actions: [
          IconButton(
            icon: const Icon(Icons.settings),
            onPressed: () {
              Navigator.push(
                context,
                MaterialPageRoute(builder: (context) => const SettingsScreen()),
              );
            },
            tooltip: 'Settings',
          ),
        ],
      ),
      body: Column(
        children: [
          // Summary card - more concise with monthly and annual costs
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: Card(
              child: Padding(
                padding: const EdgeInsets.all(16.0),
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.spaceAround,
                  children: [
                    // Monthly cost
                    Column(
                      children: [
                        const Text(
                          'Monthly',
                          style: TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const SizedBox(height: 4),
                        Text(
                          '\$${paymentProvider.totalMonthlyCost.toStringAsFixed(2)}',
                          style: const TextStyle(
                            fontSize: 20,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ],
                    ),

                    // Vertical divider
                    const SizedBox(
                      height: 50,
                      child: VerticalDivider(
                        thickness: 1,
                        width: 20,
                        color: Colors.grey,
                      ),
                    ),

                    // Annual cost
                    Column(
                      children: [
                        const Text(
                          'Annually',
                          style: TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const SizedBox(height: 4),
                        Text(
                          '\$${(paymentProvider.totalMonthlyCost * 12).toStringAsFixed(2)}',
                          style: const TextStyle(
                            fontSize: 20,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ],
                    ),

                    // Vertical divider
                    const SizedBox(
                      height: 50,
                      child: VerticalDivider(
                        thickness: 1,
                        width: 20,
                        color: Colors.grey,
                      ),
                    ),

                    // Active payments count
                    Column(
                      children: [
                        const Text(
                          'Active',
                          style: TextStyle(
                            fontSize: 16,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                        const SizedBox(height: 4),
                        Text(
                          '${paymentProvider.activePaymentsCount}',
                          style: const TextStyle(
                            fontSize: 20,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
          ),

          // Payment list
          const Expanded(
            child: PaymentList(),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          showModalBottomSheet(
            context: context,
            isScrollControlled: true,
            builder: (context) => const AddPaymentForm(),
          );
        },
        tooltip: 'Add Payment',
        child: const Icon(Icons.add),
      ),
    );
  }
}
