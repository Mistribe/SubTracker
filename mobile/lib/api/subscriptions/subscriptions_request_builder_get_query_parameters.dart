// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
/// Retrieve a paginated list of all subscriptions for the authenticated user
class SubscriptionsRequestBuilderGetQueryParameters implements AbstractQueryParameters {
    ///  Page number (default: 1)
    int? page;
    ///  Number of items per page (default: 10)
    int? size;
    /// Extracts the query parameters into a map for the URI template parsing.
    @override
    Map<String, dynamic> toMap() {
        return {
            'page' : page,
            'size' : size,
        };
    }
}
