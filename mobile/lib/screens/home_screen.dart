import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/authentication_provider.dart';
import '../pages/subscription_page.dart';
import '../pages/family_management_page.dart';
import '../pages/label_management_page.dart';
import '../widgets/app_drawer.dart';
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
      drawer: AppDrawer(
        selectedIndex: _getSelectedIndex(authProvider),
        onDestinationSelected: (index) {
          Navigator.pop(context); // Close the drawer

          // Calculate the index of the Settings destination based on drawer layout
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
