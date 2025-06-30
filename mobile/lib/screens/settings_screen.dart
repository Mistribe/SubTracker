import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/theme_provider.dart';

class SettingsScreen extends StatelessWidget {
  const SettingsScreen({super.key});

  @override
  Widget build(BuildContext context) {
    final themeProvider = Provider.of<ThemeProvider>(context);

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
          // Additional settings can be added here in the future
        ],
      ),
    );
  }
}
