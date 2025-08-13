// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
/// Retrieve a paginated list of all providers with their plans and prices
class ProvidersRequestBuilderGetQueryParameters implements AbstractQueryParameters {
    ///  Limit per request (default: 10)
    int? limit;
    ///  Offset (default: 0)
    int? offset;
    ///  Search term
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
