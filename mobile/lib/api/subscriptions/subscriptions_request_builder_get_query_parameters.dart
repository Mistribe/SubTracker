// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
/// Retrieve a paginated list of all subscriptions for the authenticated user
class SubscriptionsRequestBuilderGetQueryParameters implements AbstractQueryParameters {
    ///  Number of items per page (default: 10)
    int? limit;
    ///  Page number (default: 0)
    int? offset;
    ///  Search text
    String? search;
    /// Extracts the query parameters into a map for the URI template parsing.
    @override
    Map<String, dynamic> toMap() {
        return {
            'limit' : limit,
            'offset' : offset,
            'search' : search,
        };
    }
}
