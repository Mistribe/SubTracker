// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
/// Get all subscriptions
class SubscriptionsRequestBuilderGetQueryParameters implements AbstractQueryParameters {
    ///  Offset number
    int? page;
    ///  Number of items per page
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
