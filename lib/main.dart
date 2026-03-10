import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import 'screens/home_screen.dart';
import 'theme.dart';
import 'providers/ai_provider.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  SystemChrome.setSystemUIOverlayStyle(
    const SystemUiOverlayStyle(
      statusBarColor: Colors.transparent,
      statusBarIconBrightness: Brightness.light,
      systemNavigationBarColor: Color(0xFF151B2B),
      systemNavigationBarIconBrightness: Brightness.light,
    ),
  );
  
  // Initialize AI Provider
  final aiProvider = AIProvider();
  await aiProvider.initialize();
  
  runApp(MyApp(aiProvider: aiProvider));
}

class MyApp extends StatelessWidget {
  final AIProvider aiProvider;
  
  const MyApp({super.key, required this.aiProvider});

  @override
  Widget build(BuildContext context) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider.value(value: aiProvider),
      ],
      child: MaterialApp(
        title: 'Evolipia Radar',
        debugShowCheckedModeBanner: false,
        theme: AppTheme.darkTheme,
        home: const HomeScreen(),
      ),
    );
  }
}
