// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
/// Retrieve a paginated list of all providers with their plans and prices
class ProvidersRequestBuilderGetQueryParameters implements AbstractQueryParameters {
    ///  Page number (default: 1)
    int? page;
    ///  Items per page (default: 10)
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
