import 'package:flutter/material.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:provider/provider.dart';
import 'package:subscription_tracker/models/subscription_payment.dart';
import 'models/subscription.dart';
import 'models/settings.dart';
import 'models/label.dart';
import 'models/family_member.dart';
import 'providers/subscription_provider.dart';
import 'providers/theme_provider.dart';
import 'providers/family_member_provider.dart';
import 'providers/sync_provider.dart';
import 'repositories/subscription_repository.dart';
import 'repositories/settings_repository.dart';
import 'repositories/label_repository.dart';
import 'repositories/family_member_repository.dart';
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
  Hive.registerAdapter(LabelAdapter());
  Hive.registerAdapter(FamilyMemberAdapter());

  // Initialize repositories
  final paymentRepository = SubscriptionRepository();
  await paymentRepository.initialize();

  final settingsRepository = SettingsRepository();
  await settingsRepository.initialize();

  final labelRepository = LabelRepository();
  await labelRepository.initialize();

  final familyMemberRepository = await FamilyMemberRepository.initialize();

  runApp(
    MyApp(
      subscriptionRepository: paymentRepository,
      settingsRepository: settingsRepository,
      labelRepository: labelRepository,
      familyMemberRepository: familyMemberRepository,
    ),
  );
}

class MyApp extends StatelessWidget {
  final SubscriptionRepository subscriptionRepository;
  final SettingsRepository settingsRepository;
  final LabelRepository labelRepository;
  final FamilyMemberRepository familyMemberRepository;

  const MyApp({
    super.key,
    required this.subscriptionRepository,
    required this.settingsRepository,
    required this.labelRepository,
    required this.familyMemberRepository,
  });

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        // SyncProvider must be created before SubscriptionProvider
        ChangeNotifierProvider(
          create: (_) {
            final syncProvider = SyncProvider(
              subscriptionRepository: subscriptionRepository,
              familyMemberRepository: familyMemberRepository,
              labelRepository: labelRepository,
            );
            // Set the sync provider in the repository
            subscriptionRepository.setSyncProvider(syncProvider);
            familyMemberRepository.setSyncProvider(syncProvider);
            labelRepository.setSyncProvider(syncProvider);
            return syncProvider;
          },
        ),
        ChangeNotifierProvider(
          create: (_) => SubscriptionProvider(
            subscriptionRepository: subscriptionRepository,
            settingsRepository: settingsRepository,
            labelRepository: labelRepository,
          ),
        ),
        ChangeNotifierProvider(
          create: (_) => ThemeProvider(settingsRepository: settingsRepository),
        ),
        ChangeNotifierProvider(
          create: (_) => FamilyMemberProvider(
            familyMemberRepository: familyMemberRepository,
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
