# Subscription Tracker

A Flutter application to track subscriptions and their history.

## Features

- Add, edit, and delete subscriptions
- Track subscription history and price changes over time
- Calculate monthly and annual costs
- Persistent storage of subscription data

## Persistence Implementation

This application uses Hive for local data persistence. Hive is a lightweight, fast NoSQL database that works well with Flutter and is suitable for storing structured data like subscription records.

### Architecture

The persistence layer follows a repository pattern to abstract the data access:

1. **Models**: The `Subscription` and `PriceChange` classes in `lib/models/subscription.dart` are annotated with Hive annotations for serialization.

2. **Repository**: The `SubscriptionRepository` class in `lib/repositories/subscription_repository.dart` handles all interactions with the Hive database.

3. **Provider**: The `SubscriptionProvider` class in `lib/providers/subscription_provider.dart` uses the repository to load, save, update, and delete subscriptions.

### Future API Integration

The current architecture is designed to make it easy to switch from local storage to an API-based solution in the future. Here's how to implement API integration:

1. **Create an API Service**:
   ```dart
   class SubscriptionApiService {
     final String baseUrl;
     final http.Client client;

     SubscriptionApiService({required this.baseUrl, http.Client? client})
         : client = client ?? http.Client();

     Future<List<Subscription>> getSubscriptions() async {
       final response = await client.get(Uri.parse('$baseUrl/subscriptions'));
       if (response.statusCode == 200) {
         final List<dynamic> data = json.decode(response.body);
         return data.map((json) => Subscription.fromJson(json)).toList();
       } else {
         throw Exception('Failed to load subscriptions');
       }
     }

     Future<Subscription> addSubscription(Subscription subscription) async {
       // Implementation
     }

     Future<Subscription> updateSubscription(Subscription subscription) async {
       // Implementation
     }

     Future<void> deleteSubscription(String id) async {
       // Implementation
     }
   }
   ```

2. **Create a New Repository Implementation**:
   ```dart
   class ApiSubscriptionRepository implements SubscriptionRepository {
     final SubscriptionApiService apiService;

     ApiSubscriptionRepository({required this.apiService});

     @override
     Future<void> initialize() async {
       // No initialization needed for API
     }

     @override
     List<Subscription> getAllSubscriptions() async {
       return await apiService.getSubscriptions();
     }

     @override
     Future<void> addSubscription(Subscription subscription) async {
       await apiService.addSubscription(subscription);
     }

     @override
     Future<void> updateSubscription(Subscription subscription) async {
       await apiService.updateSubscription(subscription);
     }

     @override
     Future<void> deleteSubscription(String id) async {
       await apiService.deleteSubscription(id);
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
