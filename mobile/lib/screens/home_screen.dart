import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/authentication_provider.dart';
import '../pages/subscription_page.dart';
import '../pages/family_management_page.dart';
import '../pages/label_management_page.dart';
import '../widgets/app_drawer.dart';
import '../animations/page_transitions.dart';
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
  HomeScreenPage _previousPage = HomeScreenPage.subscriptions;

  Widget _getPageContent() {
    // Use a key based on the current page to trigger the AnimatedSwitcher
    final Widget pageContent;
    switch (_currentPage) {
      case HomeScreenPage.subscriptions:
        pageContent = const SubscriptionPage();
        break;
      case HomeScreenPage.family:
        pageContent = const FamilyManagementPage();
        break;
      case HomeScreenPage.labels:
        pageContent = const LabelManagementPage();
        break;
    }

    // Determine the slide direction based on the page index change
    final bool slideFromRight = _getPageIndex(_currentPage) > _getPageIndex(_previousPage);

    // Wrap the content in AnimatedSwitcher for smooth transitions
    return AnimatedSwitcher(
      duration: const Duration(milliseconds: 300),
      transitionBuilder: (Widget child, Animation<double> animation) {
        // Create a slide and fade transition with direction based on navigation
        return FadeTransition(
          opacity: animation,
          child: SlideTransition(
            position: Tween<Offset>(
              begin: Offset(slideFromRight ? 0.3 : -0.3, 0.0),
              end: Offset.zero,
            ).animate(CurvedAnimation(
              parent: animation,
              curve: Curves.easeOutCubic,
            )),
            child: child,
          ),
        );
      },
      child: KeyedSubtree(
        key: ValueKey<HomeScreenPage>(_currentPage),
        child: pageContent,
      ),
    );
  }

  // Helper method to get the index of a page for direction comparison
  int _getPageIndex(HomeScreenPage page) {
    switch (page) {
      case HomeScreenPage.subscriptions:
        return 0;
      case HomeScreenPage.family:
        return 1;
      case HomeScreenPage.labels:
        return 2;
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
              PageTransitions.depthSharedAxisTransition(
                page: const SettingsScreen(),
              ),
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
            // Store the current page as previous before changing
            _previousPage = _currentPage;

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
                  PageTransitions.depthSharedAxisTransition(
                    page: const SubscriptionFormScreen(),
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
