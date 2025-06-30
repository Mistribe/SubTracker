# Recurrent Payment Tracker

A Flutter application to track recurring payments and their history.

## Features

- Add, edit, and delete recurring payments
- Track payment history and price changes over time
- Calculate monthly and annual costs
- Persistent storage of payment data

## Persistence Implementation

This application uses Hive for local data persistence. Hive is a lightweight, fast NoSQL database that works well with Flutter and is suitable for storing structured data like payment records.

### Architecture

The persistence layer follows a repository pattern to abstract the data access:

1. **Models**: The `Payment` and `PriceChange` classes in `lib/models/payment.dart` are annotated with Hive annotations for serialization.

2. **Repository**: The `PaymentRepository` class in `lib/repositories/payment_repository.dart` handles all interactions with the Hive database.

3. **Provider**: The `PaymentProvider` class in `lib/providers/payment_provider.dart` uses the repository to load, save, update, and delete payments.

### Future API Integration

The current architecture is designed to make it easy to switch from local storage to an API-based solution in the future. Here's how to implement API integration:

1. **Create an API Service**:
   ```dart
   class PaymentApiService {
     final String baseUrl;
     final http.Client client;

     PaymentApiService({required this.baseUrl, http.Client? client})
         : client = client ?? http.Client();

     Future<List<Payment>> getPayments() async {
       final response = await client.get(Uri.parse('$baseUrl/payments'));
       if (response.statusCode == 200) {
         final List<dynamic> data = json.decode(response.body);
         return data.map((json) => Payment.fromJson(json)).toList();
       } else {
         throw Exception('Failed to load payments');
       }
     }

     Future<Payment> addPayment(Payment payment) async {
       // Implementation
     }

     Future<Payment> updatePayment(Payment payment) async {
       // Implementation
     }

     Future<void> deletePayment(String id) async {
       // Implementation
     }
   }
   ```

2. **Create a New Repository Implementation**:
   ```dart
   class ApiPaymentRepository implements PaymentRepository {
     final PaymentApiService apiService;

     ApiPaymentRepository({required this.apiService});

     @override
     Future<void> initialize() async {
       // No initialization needed for API
     }

     @override
     List<Payment> getAllPayments() async {
       return await apiService.getPayments();
     }

     @override
     Future<void> addPayment(Payment payment) async {
       await apiService.addPayment(payment);
     }

     @override
     Future<void> updatePayment(Payment payment) async {
       await apiService.updatePayment(payment);
     }

     @override
     Future<void> deletePayment(String id) async {
       await apiService.deletePayment(id);
     }
   }
   ```

3. **Implement Caching Strategy**:
   - Consider implementing a caching strategy to improve performance and provide offline capabilities.
   - You could use Hive to cache API responses and implement a sync mechanism.

4. **Update Dependency Injection**:
   - Modify the `main.dart` file to inject the appropriate repository implementation based on configuration.

5. **Error Handling and Connectivity**:
   - Add proper error handling for API calls
   - Implement connectivity checks to handle offline scenarios

By following this approach, you can easily switch between local storage and API-based storage without changing the rest of the application code.

## Getting Started

To run this project:

1. Ensure you have Flutter installed
2. Clone the repository
3. Run `flutter pub get` to install dependencies
4. Run `flutter run` to start the application
