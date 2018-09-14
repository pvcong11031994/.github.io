// Convert 2 byte number to 1 byte number
String.prototype.to1ByteNumber = String.prototype.to1ByteNumber || function () {
        var code0 = "０".charCodeAt(0);
        var string1byte = "";
        for (var i = 0; i < this.length; i++) {
            if (this.charCodeAt(i) >= code0) {
                string1byte = string1byte + (this.charCodeAt(i) - code0);
            } else {
                string1byte = string1byte + this[i];
            }
        }

        return string1byte;
    }

