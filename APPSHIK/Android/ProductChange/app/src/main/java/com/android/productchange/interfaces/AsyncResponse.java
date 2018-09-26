package com.android.productchange.interfaces;

import com.android.productchange.db.entity.Returnbooks;

import java.util.ArrayList;
import java.util.List;

/**
 * Created by cong-pv
 * on 2018/09/05.
 */

public interface AsyncResponse {

    void processFinish(List<Returnbooks> output);
}
