import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import '../providers/authentication_provider.dart';
import '../repositories/settings_repository.dart';
import 'home_screen.dart';

class WelcomeScreen extends StatelessWidget {
  final SettingsRepository settingsRepository;

  const WelcomeScreen({
    super.key,
    required this.settingsRepository,
  });

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 24.0),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              // App logo or icon
              Icon(
                Icons.account_balance_wallet,
                size: 80,
                color: Theme.of(context).colorScheme.primary,
              ),

              const SizedBox(height: 24),

              // App title
              const Text(
                'Subscription Tracker',
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: 28,
                  fontWeight: FontWeight.bold,
                ),
              ),

              const SizedBox(height: 16),

              // App description
              Text(
                'Track and manage all your recurring payments in one place',
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: 16,
                  color: Theme.of(context).colorScheme.onSurface.withOpacity(0.6),
                ),
              ),

              const SizedBox(height: 48),

              // Login/Signup button
              ElevatedButton(
                onPressed: () async {
                  final authProvider = Provider.of<AuthenticationProvider>(
                    context, 
                    listen: false
                  );
                  await authProvider.signIn();

                  // Mark onboarding as completed
                  await settingsRepository.completeOnboarding();

                  // Navigate to home screen if authenticated
                  if (authProvider.isAuthenticated) {
                    if (context.mounted) {
                      Navigator.of(context).pushReplacement(
                        MaterialPageRoute(
                          builder: (context) => const HomeScreen(),
                        ),
                      );
                    }
                  }
                },
                style: ElevatedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(vertical: 16),
                  backgroundColor: Theme.of(context).colorScheme.primaryContainer,
                  foregroundColor: Theme.of(context).colorScheme.onPrimaryContainer,
                ),
                child: const Text(
                  'Login / Create Account',
                  style: TextStyle(fontSize: 16),
                ),
              ),

              const SizedBox(height: 16),

              // Use anonymously button
              OutlinedButton(
                onPressed: () async {
                  // Mark onboarding as completed
                  await settingsRepository.completeOnboarding();

                  // Navigate to home screen without authentication
                  if (context.mounted) {
                    Navigator.of(context).pushReplacement(
                      MaterialPageRoute(
                        builder: (context) => const HomeScreen(),
                      ),
                    );
                  }
                },
                style: OutlinedButton.styleFrom(
                  padding: const EdgeInsets.symmetric(vertical: 16),
                  foregroundColor: Theme.of(context).colorScheme.primary,
                  side: BorderSide(color: Theme.of(context).colorScheme.outline),
                ),
                child: const Text(
                  'Continue Offline',
                  style: TextStyle(fontSize: 16),
                ),
              ),

              const SizedBox(height: 24),

              // Privacy note
              Text(
                'Your data is stored locally when using offline mode',
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontSize: 12,
                  color: Theme.of(context).colorScheme.onSurface.withOpacity(0.6),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
