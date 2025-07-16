import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../models/currency.dart';
import '../providers/theme_provider.dart';
import '../providers/subscription_provider.dart';
import '../providers/user_provider.dart';
import '../providers/sync_provider.dart';
import 'family_management_screen.dart';
import 'label_management_screen.dart';
import 'auth_screen.dart';

class SettingsScreen extends StatelessWidget {
  const SettingsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final themeProvider = Provider.of<ThemeProvider>(context);
    final subscriptionProvider = Provider.of<SubscriptionProvider>(context);

    // Get currencies from the Currency enum
    final List<String> currencies = Currency.codes;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Settings'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.of(context).pop(),
        ),
      ),
      body: ListView(
        children: [
          ListTile(
            title: const Text('Appearance'),
            leading: const Icon(Icons.palette),
          ),
          SwitchListTile(
            title: const Text('Use System Theme'),
            subtitle: const Text('Follow system dark/light mode settings'),
            secondary: const Icon(Icons.settings_system_daydream),
            value: themeProvider.useSystemTheme,
            onChanged: (value) {
              if (value) {
                themeProvider.setSystemTheme();
              } else {
                themeProvider.toggleTheme();
              }
            },
          ),
          SwitchListTile(
            title: const Text('Dark Mode'),
            subtitle: const Text('Toggle between light and dark theme'),
            secondary: Icon(
              themeProvider.isDarkMode ? Icons.dark_mode : Icons.light_mode,
            ),
            value: themeProvider.isDarkMode,
            onChanged: themeProvider.useSystemTheme ? null : (_) => themeProvider.toggleTheme(),
          ),
          const Divider(),
          ListTile(
            title: const Text('Payment Settings'),
            leading: const Icon(Icons.payments),
          ),
          ListTile(
            title: const Text('Default Currency'),
            subtitle: Text('Currency used for new subscriptions'),
            leading: const Icon(Icons.currency_exchange),
            trailing: DropdownButton<String>(
              value: subscriptionProvider.defaultCurrency,
              onChanged: (String? newValue) {
                if (newValue != null) {
                  subscriptionProvider.defaultCurrency = newValue;
                }
              },
              items: currencies.map<DropdownMenuItem<String>>((String value) {
                return DropdownMenuItem<String>(
                  value: value,
                  child: Text(value),
                );
              }).toList(),
            ),
          ),
          const Divider(),
          ListTile(
            title: const Text('Family Management'),
            leading: const Icon(Icons.family_restroom),
          ),
          ListTile(
            title: const Text('Manage Family Members'),
            subtitle: const Text('Add, edit, or remove family members'),
            leading: const Icon(Icons.people),
            trailing: const Icon(Icons.arrow_forward_ios, size: 16),
            onTap: () {
              Navigator.of(context).push(
                MaterialPageRoute(
                  builder: (context) => const FamilyManagementScreen(),
                ),
              );
            },
          ),
          const Divider(),
          ListTile(
            title: const Text('Label Management'),
            leading: const Icon(Icons.label),
          ),
          ListTile(
            title: const Text('Manage Labels'),
            subtitle: const Text('Add, edit, or remove labels'),
            leading: const Icon(Icons.label_outline),
            trailing: const Icon(Icons.arrow_forward_ios, size: 16),
            onTap: () {
              Navigator.of(context).push(
                MaterialPageRoute(
                  builder: (context) => const LabelManagementScreen(),
                ),
              );
            },
          ),
          const Divider(),
          ListTile(
            title: const Text('Account'),
            leading: const Icon(Icons.account_circle),
          ),
          Consumer<UserProvider>(
            builder: (context, userProvider, _) {
              final syncProvider = Provider.of<SyncProvider>(context);

              if (userProvider.isAuthenticated) {
                // User is signed in - show user info and sign out option
                return Column(
                  children: [
                    ListTile(
                      title: Text(userProvider.currentUser?.displayName ?? 'User'),
                      subtitle: Text(userProvider.currentUser?.email ?? ''),
                      leading: const Icon(Icons.person),
                    ),
                    ListTile(
                      title: const Text('Sign Out'),
                      subtitle: const Text('Sign out of your account'),
                      leading: const Icon(Icons.logout),
                      onTap: () async {
                        await userProvider.signOut();
                        syncProvider.updateSyncEnabled();
                      },
                    ),
                  ],
                );
              } else {
                // User is not signed in - show sign in option
                return ListTile(
                  title: const Text('Sign In / Sign Up'),
                  subtitle: const Text('Enable synchronization with your account'),
                  leading: const Icon(Icons.login),
                  trailing: const Icon(Icons.arrow_forward_ios, size: 16),
                  onTap: () {
                    Navigator.of(context).push(
                      MaterialPageRoute(
                        builder: (context) => const AuthScreen(),
                      ),
                    ).then((_) {
                      // Update sync enabled status when returning from auth screen
                      syncProvider.updateSyncEnabled();
                    });
                  },
                );
              }
            },
          ),
          const Divider(),
          // Sync status
          Consumer<SyncProvider>(
            builder: (context, syncProvider, _) {
              return ListTile(
                title: const Text('Synchronization'),
                subtitle: Text(
                  syncProvider.isSyncEnabled
                      ? 'Enabled - Data will be synced with the server'
                      : 'Disabled - Sign in to enable synchronization',
                ),
                leading: Icon(
                  syncProvider.isSyncEnabled
                      ? Icons.sync
                      : Icons.sync_disabled,
                ),
              );
            },
          ),
          // Additional settings can be added here in the future
        ],
      ),
    );
  }
}
