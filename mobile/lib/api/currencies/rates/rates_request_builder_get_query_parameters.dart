// ignore_for_file: type=lint
import 'package:microsoft_kiota_abstractions/microsoft_kiota_abstractions.dart';

/// auto generated
/// Get exchange rates for all currencies at a specific date
class RatesRequestBuilderGetQueryParameters implements AbstractQueryParameters {
    ///  Conversion date in RFC3339 format (default: current time)
    String? date;
    /// Extracts the query parameters into a map for the URI template parsing.
    @override
    Map<String, dynamic> toMap() {
        return {
            'date' : date,
        };
    }
}
