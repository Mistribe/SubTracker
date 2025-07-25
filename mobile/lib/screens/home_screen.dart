import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/authentication_provider.dart';
import '../pages/subscription_page.dart';
import '../pages/family_management_page.dart';
import '../pages/label_management_page.dart';
import 'subscription_form_screen.dart';
import 'settings_screen.dart';

enum HomeScreenPage { subscriptions, family, labels }

class HomeScreen extends StatefulWidget {
  const HomeScreen({super.key});

  @override
  State<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends State<HomeScreen> {
  HomeScreenPage _currentPage = HomeScreenPage.subscriptions;

  Widget _getPageContent() {
    switch (_currentPage) {
      case HomeScreenPage.subscriptions:
        return const SubscriptionPage();
      case HomeScreenPage.family:
        return const FamilyManagementPage();
      case HomeScreenPage.labels:
        return const LabelManagementPage();
    }
  }

  String _getPageTitle() {
    switch (_currentPage) {
      case HomeScreenPage.subscriptions:
        return 'Subscription Tracker';
      case HomeScreenPage.family:
        return 'Family Management';
      case HomeScreenPage.labels:
        return 'Manage Labels';
    }
  }

  @override
  Widget build(BuildContext context) {
    final authProvider = Provider.of<AuthenticationProvider>(context);

    return Scaffold(
      appBar: AppBar(
        title: Text(_getPageTitle()),
        backgroundColor: Theme.of(context).colorScheme.surface,
        leading: Builder(
          builder: (context) {
            return IconButton(
              icon: const Icon(Icons.menu),
              onPressed: () {
                Scaffold.of(context).openDrawer();
              },
            );
          },
        ),
      ),
      drawer: NavigationDrawer(
        selectedIndex: _getSelectedIndex(authProvider),
        onDestinationSelected: (index) {
          Navigator.pop(context); // Close the drawer

          // Calculate the index of the Settings destination
          int settingsIndex = authProvider.isAuthenticated ? 3 : 2;

          // If Settings is selected
          if (index == settingsIndex) {
            Navigator.push(
              context,
              MaterialPageRoute(builder: (context) => const SettingsScreen()),
            );
            return;
          }

          // Adjust index for authenticated vs anonymous users
          int adjustedIndex = index;
          if (!authProvider.isAuthenticated && index > 0) {
            // Skip the Family option for anonymous users
            adjustedIndex++;
          }

          setState(() {
            switch (adjustedIndex) {
              case 0: // Subscriptions
                _currentPage = HomeScreenPage.subscriptions;
                break;
              case 1: // Family
                if (authProvider.isAuthenticated) {
                  _currentPage = HomeScreenPage.family;
                }
                break;
              case 2: // Labels
                _currentPage = HomeScreenPage.labels;
                break;
            }
          });
        },
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
            selectedIcon: const Icon(Icons.subscriptions_outlined),
          ),
          // Family (only if user is connected)
          if (authProvider.isAuthenticated)
            NavigationDrawerDestination(
              icon: const Icon(Icons.family_restroom),
              label: const Text('Family'),
              selectedIcon: const Icon(Icons.family_restroom_outlined),
            ),
          // Labels
          NavigationDrawerDestination(
            icon: const Icon(Icons.label),
            label: const Text('Labels'),
            selectedIcon: const Icon(Icons.label_outline),
          ),
          const Padding(
            padding: EdgeInsets.symmetric(vertical: 8.0),
            child: Divider(),
          ),
          // Settings
          NavigationDrawerDestination(
            icon: const Icon(Icons.settings),
            label: const Text('Settings'),
            selectedIcon: const Icon(Icons.settings_outlined),
          ),
        ],
      ),
      body: SafeArea(child: _getPageContent()),
      floatingActionButton: _currentPage == HomeScreenPage.subscriptions
          ? FloatingActionButton(
              heroTag: 'home_add_subscription',
              onPressed: () {
                Navigator.push(
                  context,
                  MaterialPageRoute(
                    builder: (context) => const SubscriptionFormScreen(),
                  ),
                );
              },
              tooltip: 'Add Subscription',
              child: const Icon(Icons.add),
            )
          : null,
    );
  }

  int _getSelectedIndex(AuthenticationProvider authProvider) {
    // The Settings destination is not part of the HomeScreenPage enum,
    // so it will never be selected through this method.
    // It's handled separately in onDestinationSelected.
    switch (_currentPage) {
      case HomeScreenPage.subscriptions:
        return 0;
      case HomeScreenPage.family:
        return authProvider.isAuthenticated ? 1 : 0;
      case HomeScreenPage.labels:
        return authProvider.isAuthenticated ? 2 : 1;
    }
  }
}
