import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/family_provider.dart';
import '../providers/authentication_provider.dart';
import '../pages/subscription_page.dart';
import 'subscription_form_screen.dart';
import 'settings_screen.dart';
import 'family_management_screen.dart';
import 'label_management_screen.dart';

class HomeScreen extends StatelessWidget {
  const HomeScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final authProvider = Provider.of<AuthenticationProvider>(context);

    return Scaffold(
      appBar: AppBar(
        title: const Text(
          'Subscription Tracker',
          style: TextStyle(fontSize: 20, fontWeight: FontWeight.w600),
        ),
        elevation: 0,
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
      drawer: NavigationDrawer(
        children: [
          // User header
          Padding(
            padding: const EdgeInsets.fromLTRB(28, 16, 16, 10),
            child: Row(
              children: [
                CircleAvatar(
                  backgroundColor: Theme.of(context).colorScheme.primary,
                  child: Icon(
                    Icons.person,
                    color: Theme.of(context).colorScheme.onPrimary,
                  ),
                ),
                const SizedBox(width: 16),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        authProvider.isAuthenticated
                            ? authProvider.user?.displayName ??
                                  authProvider.user?.email ??
                                  'User'
                            : 'Anonymous',
                        style: Theme.of(context).textTheme.titleMedium,
                        overflow: TextOverflow.ellipsis,
                      ),
                      if (authProvider.isAuthenticated &&
                          authProvider.user?.email != null)
                        Text(
                          authProvider.user!.email,
                          style: Theme.of(context).textTheme.bodySmall,
                          overflow: TextOverflow.ellipsis,
                        ),
                    ],
                  ),
                ),
              ],
            ),
          ),
          const Divider(),
          // Subscriptions (Home)
          NavigationDrawerDestination(
            icon: const Icon(Icons.subscriptions),
            label: const Text('Subscriptions'),
          ),
          // Family (only if user is connected)
          if (authProvider.isAuthenticated)
            NavigationDrawerDestination(
              icon: const Icon(Icons.family_restroom),
              label: const Text('Family'),
            ),
          // Labels
          NavigationDrawerDestination(
            icon: const Icon(Icons.label),
            label: const Text('Labels'),
          ),
        ],
        onDestinationSelected: (index) {
          Navigator.pop(context); // Close the drawer

          // Adjust index for authenticated vs anonymous users
          int adjustedIndex = index;
          if (!authProvider.isAuthenticated && index > 0) {
            // Skip the Family option for anonymous users
            adjustedIndex++;
          }

          switch (adjustedIndex) {
            case 0: // Subscriptions (Home)
              // Already on home screen
              break;
            case 1: // Family
              if (authProvider.isAuthenticated) {
                Navigator.push(
                  context,
                  MaterialPageRoute(
                    builder: (context) => const FamilyManagementScreen(),
                  ),
                );
              }
              break;
            case 2: // Labels
              Navigator.push(
                context,
                MaterialPageRoute(
                  builder: (context) => const LabelManagementScreen(),
                ),
              );
              break;
          }
        },
      ),
      body: SafeArea(child: const SubscriptionPage()),
      floatingActionButton: FloatingActionButton(
        onPressed: () {
          Navigator.push(
            context,
            MaterialPageRoute(
              builder: (context) => const SubscriptionFormScreen(),
            ),
          );
        },
        tooltip: 'Add Subscription',
        elevation: 4,
        child: const Icon(Icons.add),
      ),
    );
  }
}
