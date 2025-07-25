import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/authentication_provider.dart';
import '../screens/settings_screen.dart';

class AppDrawer extends StatelessWidget {
  final int selectedIndex;
  final Function(int) onDestinationSelected;

  const AppDrawer({
    super.key,
    required this.selectedIndex,
    required this.onDestinationSelected,
  });

  @override
  Widget build(BuildContext context) {
    final authProvider = Provider.of<AuthenticationProvider>(context);

    // Calculate settings index for proper selection highlighting
    int adjustedSelectedIndex = selectedIndex;
    int settingsIndex = authProvider.isAuthenticated ? 3 : 2;

    return Stack(
      children: [
        // Main Navigation Drawer
        NavigationDrawer(
          selectedIndex: adjustedSelectedIndex,
          onDestinationSelected: onDestinationSelected,
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
            // Extra space to ensure settings appears at the bottom
            SizedBox(height: MediaQuery.of(context).size.height - 350),
            const Padding(
              padding: EdgeInsets.symmetric(vertical: 8.0),
              child: Divider(),
            ),
            // Settings at the bottom
            NavigationDrawerDestination(
              icon: const Icon(Icons.settings),
              label: const Text('Settings'),
              selectedIcon: const Icon(Icons.settings_outlined),
            ),
            // Bottom padding
            const SizedBox(height: 16),
          ],
        ),
      ],
    );
  }
}
