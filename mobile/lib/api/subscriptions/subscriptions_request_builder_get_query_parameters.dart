// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
/// Retrieve a paginated list of all subscriptions for the authenticated user
class SubscriptionsRequestBuilderGetQueryParameters implements AbstractQueryParameters {
    ///  Filter by start date (RFC3339)
    /// @QueryParameter('from_date')
    String? fromDate;
    ///  Number of items per page (default: 10)
    int? limit;
    ///  Page number (default: 0)
    int? offset;
    ///  Filter by provider IDs
    List<String>? providers;
    ///  Filter by recurrency types
    List<String>? recurrencies;
    ///  Search text
    String? search;
    ///  Filter by end date (RFC3339)
    /// @QueryParameter('to_date')
    String? toDate;
    ///  Filter by user IDs
    List<String>? users;
    ///  Include inactive subscriptions
    /// @QueryParameter('with_inactive')
    bool? withInactive;
    /// Extracts the query parameters into a map for the URI template parsing.
    @override
    Map<String, dynamic> toMap() {
        return {
            'from_date' : fromDate,
            'limit' : limit,
            'offset' : offset,
            'providers' : providers,
            'recurrencies' : recurrencies,
            'search' : search,
            'to_date' : toDate,
            'users' : users,
            'with_inactive' : withInactive,
        };
    }
}
