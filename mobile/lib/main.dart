import 'package:flutter/material.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:provider/provider.dart';
import 'models/payment.dart';
import 'providers/payment_provider.dart';
import 'providers/theme_provider.dart';
import 'repositories/payment_repository.dart';
import 'screens/home_screen.dart';

void main() async {
  // Ensure Flutter is initialized
  WidgetsFlutterBinding.ensureInitialized();

  // Initialize Hive
  await Hive.initFlutter();

  // Register Hive adapters
  Hive.registerAdapter(PaymentAdapter());
  Hive.registerAdapter(PaymentDetailAdapter());

  // Initialize repository
  final paymentRepository = PaymentRepository();
  await paymentRepository.initialize();

  runApp(MyApp(paymentRepository: paymentRepository));
}

class MyApp extends StatelessWidget {
  final PaymentRepository paymentRepository;

  const MyApp({super.key, required this.paymentRepository});

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider(
          create: (_) => PaymentProvider(paymentRepository: paymentRepository),
        ),
        ChangeNotifierProvider(create: (_) => ThemeProvider()),
      ],
      child: Consumer<ThemeProvider>(
        builder: (context, themeProvider, _) {
          return MaterialApp(
            title: 'Recurrent Payment Tracker',
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
