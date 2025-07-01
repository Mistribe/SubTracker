import 'package:flutter/material.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:provider/provider.dart';
import 'package:subscription_tracker/models/subscription_payment.dart';
import 'models/subscription.dart';
import 'models/settings.dart';
import 'providers/subscription_provider.dart';
import 'providers/theme_provider.dart';
import 'repositories/subscription_repository.dart';
import 'repositories/settings_repository.dart';
import 'screens/home_screen.dart';

void main() async {
  // Ensure Flutter is initialized
  WidgetsFlutterBinding.ensureInitialized();

  // Initialize Hive
  await Hive.initFlutter();

  // Register Hive adapters
  Hive.registerAdapter(SubscriptionAdapter());
  Hive.registerAdapter(SubscriptionPaymentAdapter());
  Hive.registerAdapter(SettingsAdapter());

  // Initialize repositories
  final paymentRepository = SubscriptionRepository();
  await paymentRepository.initialize();

  final settingsRepository = SettingsRepository();
  await settingsRepository.initialize();

  runApp(MyApp(
    subscriptionRepository: paymentRepository,
    settingsRepository: settingsRepository,
  ));
}

class MyApp extends StatelessWidget {
  final SubscriptionRepository subscriptionRepository;
  final SettingsRepository settingsRepository;

  const MyApp({
    super.key, 
    required this.subscriptionRepository,
    required this.settingsRepository,
  });

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(
          create: (_) => SubscriptionProvider(
            subscriptionRepository: subscriptionRepository,
            settingsRepository: settingsRepository,
          ),
        ),
        ChangeNotifierProvider(
          create: (_) => ThemeProvider(
            settingsRepository: settingsRepository,
          ),
        ),
      ],
      child: Consumer<ThemeProvider>(
        builder: (context, themeProvider, _) {
          return MaterialApp(
            title: 'Subscription Tracker',
            debugShowCheckedModeBanner: false,
            themeMode: themeProvider.themeMode,
            theme: themeProvider.lightTheme,
            darkTheme: themeProvider.darkTheme,
            home: const HomeScreen(),
          );
        },
      ),
    );
  }
}
